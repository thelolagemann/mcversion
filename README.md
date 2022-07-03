# mcversion

![GitHub release (latest SemVer including pre-releases)](https://img.shields.io/github/v/release/thelolagemann/mcversion?include_prereleases&style=for-the-badge)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://pkg.go.dev/github.com/thelolagemann/mcversion)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/thelolagemann/mcversion?style=for-the-badge)
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
	version, err := mcversion.Version("1.12.2")
	if err != nil {
		panic(err)
	}
	fmt.Println(version.Name)
}

```

	// Output: 1.12.2