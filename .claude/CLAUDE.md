# go-tags

Generic key/value pair conversion. Converts between `key=value` arguments, `map[string]string`, and lists of any type exposing `GetKey()`/`GetValue()` methods (the accessor shape generated protobuf messages share) — the conversion consumers like `skytl` and `skykerneld` would otherwise reimplement per tag command.

- Package `tags` (import `github.com/gomatic/go-tags`): `Pair` constraint (`GetKey`/`GetValue`), `ToMap[T Pair]([]T) map`, `FromMap[T any](map, func(key, value string) T) []T` (key-sorted, deterministic; the constructor is injected because a method constraint can read but not build a value), `Merge[T Pair](map, []T) map` (updates win, non-nil, base untouched), `Parse([]string) (map, error)` (key=value, later duplicate wins).
- **Generic by design — never import a schema.** This is a gomatic basic library: no `skykernel/*` (or any producer-specific) dependency may be added. Producer-coupled glue (e.g. a `typev1.Tag` constructor) belongs in the consumer's own `internal/apitags` adapter, not here.
- **Owns only the pair-list conversions.** Pure map ops — clone, key removal, map+map merge — stay in stdlib `maps`/`slices` at the call site (`maps.Clone`, `maps.Copy`); do not add `Clone`/`Without` here.
- Sole dependency is [gomatic/go-error](https://github.com/gomatic/go-error): `ErrInvalidPair` is an `errs.Const`, matched with `errors.Is`; never `errors.New` or `fmt.Errorf` — wrapping goes through `ErrInvalidPair.With(cause, args...)` (`errs.Const.With`).
- Value-oriented, immutable, private by default. Gate: gofumpt, vet, staticcheck, govulncheck, gocognit ≤ 7, 100% coverage.
- `Makefile`, `.golangci.yaml`, `.editorconfig`, `.gitignore`, `.github/` are owned and pushed by `nicerobot/tools.repository` — do not edit in-tree; per-repo divergence goes in a `Makefile.local`.
