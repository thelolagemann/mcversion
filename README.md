# mcversion

![GitHub tag (latest SemVer pre-release)](https://img.shields.io/github/v/tag/thelolagemann/mcversion?include_prereleases&label=release&style=for-the-badge)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://pkg.go.dev/github.com/thelolagemann/mcversion)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/thelolagemann/mcversion?style=for-the-badge)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/thelolagemann/mcversion/Test?label=tests&style=for-the-badge)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/thelolagemann/mcversion/CodeQL?label=CodeQL&style=for-the-badge)
[![Go ReportCard](https://goreportcard.com/badge/github.com/thelolagemann/mcversion?style=for-the-badge)](https://goreportcard.com/report/thelolagemann/mcversion)

> A small golang library to access the minecraft version API.

## Getting Started

### Installation

```go
go get github.com/mcversion/mcversion
```

### Usage

```go
package main
import (
	"fmt"
	"github.com/thelolagemann/mcversion"
)

func main() {
	manifest, _ := mcversion.Manifest()
	
	for _, version := range manifest.Versions {
		fmt.Println(version.Id)
	}
}

```

	// Output: 1.12.2