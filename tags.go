// Package tags converts between key=value arguments, string maps, and the
// skykernel API's repeated Tag lists — the small conversion that every tool and
// service speaking the skykernel API would otherwise reimplement. Pure map
// operations (clone, merge, key removal) are intentionally left to the standard
// library's maps and slices packages; this package owns only the proto-coupled
// conversions those packages cannot express.
package tags

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	errs "github.com/gomatic/go-error"
	typev1 "github.com/skykernel/api/src/proto/skykernel/type/v1"
)

// ErrInvalidPair is returned by Parse for an argument that is not "key=value"
// (missing "=" or with an empty key).
const ErrInvalidPair errs.Const = "tags: argument is not key=value"

// ToMap collapses a Tag list into a key/value map; on duplicate keys the later
// entry wins. The result is non-nil even when tags is empty.
func ToMap(tags []*typev1.Tag) map[string]string {
	out := make(map[string]string, len(tags))
	for _, t := range tags {
		out[t.GetKey()] = t.GetValue()
	}
	return out
}

// ToProto renders a map as a Tag list ordered by key, so equal maps always
// produce an identical list (deterministic output).
func ToProto(m map[string]string) []*typev1.Tag {
	out := make([]*typev1.Tag, 0, len(m))
	for _, key := range slices.Sorted(maps.Keys(m)) {
		out = append(out, &typev1.Tag{Key: key, Value: m[key]})
	}
	return out
}

// Merge returns base with updates applied on top (on a key collision the update
// wins), as an independent non-nil map; base is never modified. It is the
// proto-coupled combine that stdlib maps.Copy cannot express directly, since the
// updates arrive as a Tag list and base may be nil.
func Merge(base map[string]string, updates []*typev1.Tag) map[string]string {
	out := maps.Clone(base)
	if out == nil {
		out = make(map[string]string, len(updates))
	}
	for _, t := range updates {
		out[t.GetKey()] = t.GetValue()
	}
	return out
}

// Parse turns "key=value" arguments into a map; on duplicate keys the later
// entry wins. ErrInvalidPair (wrapping the offending argument) is returned for
// an argument missing "=" or with an empty key. The result is non-nil even when
// args is empty.
func Parse(args []string) (map[string]string, error) {
	out := make(map[string]string, len(args))
	for _, arg := range args {
		key, value, ok := strings.Cut(arg, "=")
		if !ok || key == "" {
			return nil, fmt.Errorf("%w: %q", ErrInvalidPair, arg)
		}
		out[key] = value
	}
	return out, nil
}
