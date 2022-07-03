package mcversion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func generateTestData() error {
	// create dirs
	if err := os.MkdirAll("testdata", 0755); err != nil {
		return err
	}

	// start generating files
	manifest, err := Manifest()
	if err != nil {
		return err
	}
	b, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	if err := os.WriteFile("testdata/version_manifest.json", b, 0644); err != nil {
		return err
	}

	manifestV2, err := ManifestV2()
	if err != nil {
		return err
	}
	b, err = json.Marshal(manifestV2)
	if err != nil {
		return err
	}
	if err := os.WriteFile("testdata/version_manifest_v2.json", b, 0644); err != nil {
		return err
	}
	versions, err := manifest.AllVersions()
	if err != nil {
		return err
	}
	for _, v := range versions {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join("testdata", v.ID+".json"), b, 0644); err != nil {
			return err
		}
	}

	return nil
}

type mockClient struct {
	responses map[string][]byte
	response  *http.Response
	err       error
}

func (m *mockClient) Get(url string) (*http.Response, error) {
	if m.response != nil {
		return m.response, nil
	}
	if m.err != nil {
		return nil, m.err
	}
	dec, err := url2.QueryUnescape(url)
	if err != nil {
		return nil, err
	}
	if b, ok := m.responses[filepath.Base(dec)]; ok {
		return &http.Response{Body: ioutil.NopCloser(bytes.NewReader(b)), Header: map[string][]string{"Content-Type": {"application/json"}}}, nil
	}
	return nil, fmt.Errorf("no response for %s: %s", dec, url)
}

func (m *mockClient) loadResponses() error {
	return filepath.WalkDir("testdata", func(path string, dir os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if dir.IsDir() {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		m.responses[dir.Name()] = b
		return nil
	})
}

func (m *mockClient) doGet(res *http.Response, err error, f func() error) error {
	m.response = res
	m.err = err
	defer func() { m.response = nil; m.err = nil }()
	return f()
}

var (
	clientMock = &mockClient{
		responses: map[string][]byte{},
	}
)

func init() {
	_, err := os.Open(filepath.Join("testdata", "version_manifest.json")) // TODO ensure all files are present
	if err != nil {
		if err := generateTestData(); err != nil {
			panic(err)
		}
	}
	httpClient = clientMock
	if err := clientMock.loadResponses(); err != nil {
		panic(err)
	}
}

func Test_MCVersion(t *testing.T) {
	t.Run("getJSON", func(t *testing.T) {
		t.Run("client_error", func(t *testing.T) {
			testErr := fmt.Errorf("test error")
			err := clientMock.doGet(nil, testErr, func() error {
				_, err := Manifest()
				return err
			})
			if err != testErr {
				t.Errorf("expected error, got %v", err)
			}
		})
		t.Run("network_error", func(t *testing.T) {
			err := clientMock.doGet(&http.Response{StatusCode: 500}, nil, func() error {
				_, err := Manifest()
				return err
			})
			if err == nil {
				t.Errorf("expected error")
			}
		})
		t.Run("incorrect_header", func(t *testing.T) {
			err := clientMock.doGet(&http.Response{Header: map[string][]string{"Content-Type": {"text/plain"}}}, nil, func() error {
				_, err := Manifest()
				return err
			})
			if err == nil {
				t.Errorf("expected error")
			}
		})
		t.Run("decode_error", func(t *testing.T) {
			err := clientMock.doGet(&http.Response{
				Header: map[string][]string{"Content-Type": {"application/json"}},
				Body:   io.NopCloser(bytes.NewReader([]byte{})),
			}, nil, func() error {
				_, err := Manifest()
				return err
			})
			if err == nil {
				t.Errorf("expected error")
			}
		})
	})
	t.Run("Manifest", func(t *testing.T) {
		m, err := Manifest()
		if err != nil {
			t.Fatalf("failed to load manifest: %v", err)
		}
		if len(m.Versions) == 0 {
			t.Error("no versions found")
		}
		t.Run("ManifestInfo", func(t *testing.T) {
			for _, v := range m.Versions {
				_, err := v.Info()
				if err != nil {
					t.Errorf("failed to load info for %s: %v", v.Id, err)
				}
			}
		})
		t.Run("Error", func(t *testing.T) {
			gManifestErr = fmt.Errorf("test error")
			if _, err := m.AllVersions(); err == nil {
				t.Error("expected error")
			}
			gManifestErr = nil
		})
	})
	t.Run("ManifestV2", func(t *testing.T) {
		m, err := ManifestV2()
		if err != nil {
			t.Fatalf("failed to load manifest: %v", err)
		}
		if len(m.Versions) == 0 {
			t.Error("no versions found")
		}
	})
	t.Run("AllVersions", func(t *testing.T) {
		versions, err := AllVersions()
		if err != nil {
			t.Fatal(err)
		}
		if len(versions) != len(gManifest.Versions) {
			t.Errorf("expected %d versions, got %d", len(gManifest.Versions), len(versions))
		}
		t.Run("PrematureError", func(t *testing.T) {
			go func() {
				time.Sleep(time.Millisecond)
				gManifestErr = fmt.Errorf("test error")
			}()
			_, err := AllVersions()
			if err == nil {
				t.Error("expected error")
			}
			gManifestErr = nil
		})
		t.Run("VersionInfo", func(t *testing.T) {
			for _, v := range versions {
				vi, err := Version(v.ID)
				if err != nil {
					t.Error(err)
				}
				if vi.ID != v.ID {
					t.Errorf("expected %s, got %s", v.ID, vi.ID)
				}
			}
		})
		t.Run("VersionInfo_Invalid", func(t *testing.T) {
			_, err := Version("invalid")
			if err == nil {
				t.Error("expected error")
			}
		})
	})
	t.Run("LatestRelease", func(t *testing.T) {
		version, err := LatestRelease()
		if err != nil {
			t.Error(err)
		}
		if version.Type != "release" {
			t.Errorf("expected release, got %s", version.Type)
		}

	})
	t.Run("LatestSnapshot", func(t *testing.T) {
		version, err := LatestSnapshot()
		if err != nil {
			t.Fatal(err)
		}
		if version.Type != "snapshot" {
			t.Errorf("expected snapshot, got %s", version.Type)
		}
	})

}

func Benchmark_MCVersion(b *testing.B) {
	b.Run("Manifest", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := Manifest()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("ManifestV2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := ManifestV2()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("AllVersions", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := AllVersions()
			if err != nil {
				b.Error(err)
			}
		}

	})
	b.Run("VersionInfo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := Version("1.19")
			if err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("LatestRelease", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := LatestRelease()
			if err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("LatestSnapshot", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := LatestSnapshot()
			if err != nil {
				b.Error(err)
			}
		}
	})
}
