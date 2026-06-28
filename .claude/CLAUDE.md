# go-tags

Tag conversion for the [skykernel API](https://github.com/skykernel/api). Converts between `key=value` arguments, `map[string]string`, and the API's `[]*typev1.Tag` lists — the conversion `skytl` and `skykerneld` would otherwise reimplement per tag command.

- Package `tags` (import `github.com/skykernel/go-tags`): `ToMap([]*typev1.Tag) map`, `ToProto(map) []*typev1.Tag` (key-sorted, deterministic), `Parse([]string) (map, error)` (key=value, later duplicate wins).
- **Owns only the proto-coupled conversions.** Pure map ops (clone, merge, key removal) stay in stdlib `maps`/`slices` at the call site — do not add `Clone`/`Merge`/`Without` here.
- Sole dependency is `github.com/skykernel/api` (same org). Errors are a local `type Error string` sentinel (`ErrInvalidPair`), matched with `errors.Is`; never `errors.New`, `fmt.Errorf` only to `%w`-wrap. No `gomatic/go-error` dependency — keeps the module single-org for CI fetch.
- Value-oriented, immutable, private by default. Gate: gofumpt, vet, staticcheck, govulncheck, gocognit ≤ 7, 100% coverage.
- `Makefile`, `.golangci.yaml`, `.editorconfig`, `.gitignore`, `.github/` are owned and pushed by `nicerobot/tools.repository` — do not edit in-tree; per-repo divergence goes in a `Makefile.local`.
