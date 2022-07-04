# mcversion

![GitHub tag (latest SemVer pre-release)](https://img.shields.io/github/v/tag/thelolagemann/mcversion?include_prereleases&label=release&sort=semver&style=for-the-badge)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://pkg.go.dev/github.com/thelolagemann/mcversion)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/thelolagemann/mcversion?style=for-the-badge)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/thelolagemann/mcversion/Test?label=tests&style=for-the-badge)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/thelolagemann/mcversion/CodeQL?label=CodeQL&style=for-the-badge)
[![Go ReportCard](https://goreportcard.com/badge/github.com/thelolagemann/mcversion?style=for-the-badge)](https://goreportcard.com/report/thelolagemann/mcversion)

> A small golang library to access the minecraft versions manifest, both the `version_manifest` and `version_manifest_v2`.

## Getting Started

### Installation

```shell
go get github.com/mcversion/mcversion
```

### Usage


```go
// Get the versions manifest
manifest, _ := mcversion.Manifest() // or mcversion.ManifestV2()
fmt.Printf("%+v\n", manifest.Versions[len(manifest.Versions)-1])

// Output: {Id:rd-132211 Type:old_alpha URL:https://launchermeta.mojang.com/v1/packages/d090f5d3766a28425316473d9ab6c37234d48b02/rd-132211.json Time:2022-03-10 09:51:38 +0000 GMT ReleaseTime:2009-05-13 20:11:00 +0000 +0000}

//  Get all releases
var releases []string
for _, v := range manifest.Versions {
	if v.Type == "release" {
		releases = append(releases, v.Id)
	}
}
fmt.Println(strings.Join(releases, ","))

// Output: 1.19,1.18.2,1.18.1,1.18,1.17.1,1.17,1.16.5,1.16.4,1.16.3,1.16.2,1.16.1,1.16,1.15.2,1.15.1,1.15,1.14.4,1.14.3,1.14.2,1.14.1,1.14,1.13.2,1.13.1,1.13,1.12.2,1.12.1,1.12,1.11.2,1.11.1,1.11,1.10.2,1.10.1,1.10,1.9.4,1.9.3,1.9.2,1.9.1,1.9,1.8.9,1.8.8,1.8.7,1.8.6,1.8.5,1.8.4,1.8.3,1.8.2,1.8.1,1.8,1.7.10,1.7.9,1.7.8,1.7.7,1.7.6,1.7.5,1.7.4,1.7.3,1.7.2,1.6.4,1.6.2,1.6.1,1.5.2,1.5.1,1.4.7,1.4.6,1.4.5,1.4.4,1.4.2,1.3.2,1.3.1,1.2.5,1.2.4,1.2.3,1.2.2,1.2.1,1.1,1.0
	
// Get extended details of a specific version
version, _ := mcversion.Version("1.16")
fmt.Printf("%+v\n", version.Downloads)

// or call the Info method on a ManifestVersion or ManifestVersionV2 struct
for _, v := range manifest.Versions {
	if v.Id == "1.16" {
		version, _ := v.Info()
		fmt.Printf("%+v\n", version.Downloads)
	}
}

// Output: {Client:{Sha1:c9abbe8ee4fa490751ca70635340b7cf00db83ff Size:17492432 URL:https://launcher.mojang.com/v1/objects/c9abbe8ee4fa490751ca70635340b7cf00db83ff/client.jar} ClientMappings:{Sha1:ddf517a4f6750f4c15189de4e03246ae1f916cf5 Size:5632455 URL:https://launcher.mojang.com/v1/objects/ddf517a4f6750f4c15189de4e03246ae1f916cf5/client.txt} Server:{Sha1:a412fd69db1f81db3f511c1463fd304675244077 Size:37968964 URL:https://launcher.mojang.com/v1/objects/a412fd69db1f81db3f511c1463fd304675244077/server.jar} ServerMappings:{Sha1:11120c39da4df293c4bd020896391fb9ddd6c2ba Size:4329615 URL:https://launcher.mojang.com/v1/objects/11120c39da4df293c4bd020896391fb9ddd6c2ba/server.txt}}
```
