# go-tags

Tag conversion for the [skykernel API](https://github.com/skykernel/api): turn `key=value` arguments and string maps into the API's repeated `Tag` lists, and back. It is the small conversion every skykernel tool (`skytl`) and service (`skykerneld`) would otherwise reimplement per command.

Pure map operations — clone, merge, key removal — are deliberately **not** here: the standard library's [`maps`](https://pkg.go.dev/maps) and [`slices`](https://pkg.go.dev/slices) already express them (`maps.Clone`, `maps.Copy`, `slices.Sorted(maps.Keys(m))`). This package owns only the proto-coupled conversions those packages cannot.

## Install

```sh
go get github.com/skykernel/go-tags
```

Requires Go 1.26+.

## Usage

```go
import (
	typev1 "github.com/skykernel/api/src/proto/skykernel/type/v1"

	"github.com/skykernel/go-tags"
)

// CLI args -> map (later duplicate wins; ErrInvalidPair on a bad arg).
m, err := tags.Parse([]string{"env=prod", "team=core"})

// map -> deterministic, key-sorted Tag list (for a request).
req := tags.ToProto(m)

// Tag list (from a response) -> map (for output).
out := tags.ToMap(resp.GetTagSet())
```

Errors are package sentinels matchable with `errors.Is`:

```go
if _, err := tags.Parse(args); errors.Is(err, tags.ErrInvalidPair) {
	// one of the args was not key=value
}
```
