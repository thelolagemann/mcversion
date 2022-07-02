# mcversion

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