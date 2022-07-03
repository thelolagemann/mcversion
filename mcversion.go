package mcversion

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const (
	manifestURL   = "https://launchermeta.mojang.com/mc/game/version_manifest.json"
	manifestV2URL = "https://launchermeta.mojang.com/mc/game/version_manifest_v2.json"
)

type client interface {
	Get(url string) (*http.Response, error)
}

var (
	httpClient client = &http.Client{
		Timeout: time.Second * 5,
	}

	gManifest    = VersionManifest{}
	isLoaded     bool
	gManifestErr error
	versionPool  = &sync.Pool{
		New: func() interface{} {
			return VersionInfo{}
		},
	}
)

// VersionManifest is the manifest of all versions.
type VersionManifest struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`
	Versions []ManifestVersion `json:"versions"`
}

// ManifestVersion is a version in the manifest. In order to
// get the VersionInfo, which contains extended information about
// a specific version, call the Info method on the ManifestVersion,
// or alternatively call the Version method with the specified ID.
type ManifestVersion struct {
	Id          string    `json:"id"`
	Type        string    `json:"type"`
	URL         string    `json:"url"`
	Time        time.Time `json:"time"`
	ReleaseTime time.Time `json:"releaseTime"`
}

// Info returns the VersionInfo for the calling ManifestVersion. See
// also: Version.
func (m ManifestVersion) Info() (VersionInfo, error) {
	return Version(m.Id)
}

// VersionManifestV2 like VersionManifest, but includes the SHA1, and
// compliance level in the Versions field.
type VersionManifestV2 struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`
	Versions []ManifestVersionV2 `json:"versions"`
}

// ManifestVersionV2 like ManifestVersion, with the addition of
// the SHA1 and compliance level fields.
type ManifestVersionV2 struct {
	ManifestVersion
	SHA1            []byte `json:"sha1"`
	ComplianceLevel int    `json:"complianceLevel"`
}

// LatestRelease returns the latest release.
func (vM VersionManifest) LatestRelease() (VersionInfo, error) {
	return vM.Version(vM.Latest.Release)
}

// LatestSnapshot returns the latest snapshot.
func (vM VersionManifest) LatestSnapshot() (VersionInfo, error) {
	return vM.Version(vM.Latest.Snapshot)
}

// Version fetches a version by id.
func (vM VersionManifest) Version(id string) (VersionInfo, error) {
	version := versionPool.Get().(VersionInfo)
	for _, v := range vM.Versions {
		if v.Id == id {
			if err := vM.getJSON(v.URL, &version); err != nil {
				return version, err
			}
			versionPool.Put(version)
			return version, nil
		}
	}
	versionPool.Put(version)
	return version, fmt.Errorf("version %s not found", id)
}

// AllVersions returns all versions.
func (vM VersionManifest) AllVersions() ([]VersionInfo, error) {
	errs, ctx := errgroup.WithContext(context.Background())
	errs.SetLimit(runtime.NumCPU())
	results := make(chan VersionInfo, len(vM.Versions))

	// iterate over manifest versions
	for _, v := range vM.Versions {
		id := v.Id
		errs.Go(func() error {
			version, err := vM.Version(id)
			if err != nil {
				return err
			}
			select {
			case results <- version:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		})
	}

	// close out the channel if any error occurs
	go func() {
		errs.Wait()
		close(results)
	}()

	versions := make([]VersionInfo, 0, len(vM.Versions))
	// wait for all fetches to finish
	for r := range results {
		versions = append(versions, r)
	}
	return versions, errs.Wait()
}

// VersionInfo contains extended information about a version,
// not present in the VersionManifest.
type VersionInfo struct {
	Arguments struct {
		Game []interface{} `json:"game"`
		Jvm  []interface{} `json:"jvm"`
	} `json:"arguments"`
	AssetIndex struct {
		ID        string `json:"id"`
		Sha1      string `json:"sha1"`
		Size      int    `json:"size"`
		TotalSize int    `json:"totalSize"`
		URL       string `json:"url"`
	} `json:"assetIndex"`
	Assets          string `json:"assets"`
	ComplianceLevel int    `json:"complianceLevel"`
	Downloads       struct {
		Client struct {
			Sha1 string `json:"sha1"`
			Size int    `json:"size"`
			URL  string `json:"url"`
		} `json:"client"`
		ClientMappings struct {
			Sha1 string `json:"sha1"`
			Size int    `json:"size"`
			URL  string `json:"url"`
		} `json:"client_mappings"`
		Server struct {
			Sha1 string `json:"sha1"`
			Size int    `json:"size"`
			URL  string `json:"url"`
		} `json:"server"`
		ServerMappings struct {
			Sha1 string `json:"sha1"`
			Size int    `json:"size"`
			URL  string `json:"url"`
		} `json:"server_mappings"`
	} `json:"downloads"`
	ID          string `json:"id"`
	JavaVersion struct {
		Component    string `json:"component"`
		MajorVersion int    `json:"majorVersion"`
	} `json:"javaVersion"`
	Libraries []struct {
		Downloads struct {
			Artifact struct {
				Path string `json:"path"`
				Sha1 string `json:"sha1"`
				Size int    `json:"size"`
				URL  string `json:"url"`
			} `json:"artifact"`
		} `json:"downloads"`
		Name  string `json:"name"`
		Rules []struct {
			Action string `json:"action"`
			Os     struct {
				Name string `json:"name"`
			} `json:"os"`
		} `json:"rules,omitempty"`
		Natives struct {
			Osx string `json:"osx"`
		} `json:"natives,omitempty"`
		Extract struct {
			Exclude []string `json:"exclude"`
		} `json:"extract,omitempty"`
	} `json:"libraries"`
	Logging struct {
		Client struct {
			Argument string `json:"argument"`
			File     struct {
				ID   string `json:"id"`
				Sha1 string `json:"sha1"`
				Size int    `json:"size"`
				URL  string `json:"url"`
			} `json:"file"`
			Type string `json:"type"`
		} `json:"client"`
	} `json:"logging"`
	MainClass              string    `json:"mainClass"`
	MinimumLauncherVersion int       `json:"minimumLauncherVersion"`
	ReleaseTime            time.Time `json:"releaseTime"`
	Time                   time.Time `json:"time"`
	Type                   string    `json:"type"`
}

// AllVersions returns all versions.
func AllVersions() ([]VersionInfo, error) {
	return globalManifest().AllVersions()
}

// LatestRelease returns the latest release.
func LatestRelease() (VersionInfo, error) {
	return globalManifest().LatestRelease()
}

// LatestSnapshot returns the latest snapshot.
func LatestSnapshot() (VersionInfo, error) {
	return globalManifest().LatestSnapshot()
}

// Version returns a version by id.
func Version(id string) (VersionInfo, error) {
	return globalManifest().Version(id)
}

// Manifest returns the manifest of all versions.
func Manifest() (VersionManifest, error) {
	var versionManifest VersionManifest
	err := getJSON(manifestURL, &versionManifest)
	return versionManifest, err
}

// ManifestV2 like Manifest, but includes the SHA1, and
// compliance level in the fields of each ManifestVersionV2.
func ManifestV2() (VersionManifestV2, error) {
	var versionManifest VersionManifestV2
	err := getJSON(manifestV2URL, &versionManifest)
	return versionManifest, err
}

func (vM VersionManifest) getJSON(url string, v interface{}) error {
	if gManifestErr != nil {
		return gManifestErr
	}
	return getJSON(url, v)
}

func getJSON(url string, v interface{}) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("%s: %s", url, resp.Status)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("%s: unexpected content type: %s", url, resp.Header.Get("Content-Type"))
	}
	if err = json.NewDecoder(resp.Body).Decode(v); err != nil {
		return err
	}
	return resp.Body.Close()
}

// globalManifest is a small helper function that
// returns the global manifest if it has been loaded,
// otherwise it loads it, and then returns the manifest.
func globalManifest() VersionManifest {
	if !isLoaded {
		var err error
		gManifest, err = Manifest()
		gManifestErr = err
		isLoaded = true
		return gManifest
	}
	return gManifest
}
