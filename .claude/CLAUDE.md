# go-tags

Tag conversion for the [skykernel API](https://github.com/skykernel/api). Converts between `key=value` arguments, `map[string]string`, and the API's `[]*typev1.Tag` lists — the conversion `skytl` and `skykerneld` would otherwise reimplement per tag command.

- Package `tags` (import `github.com/gomatic/go-tags`): `ToMap([]*typev1.Tag) map`, `ToProto(map) []*typev1.Tag` (key-sorted, deterministic), `Merge(map, []*typev1.Tag) map` (updates win, non-nil, base untouched), `Parse([]string) (map, error)` (key=value, later duplicate wins).
- **Owns only the proto-coupled conversions** (those taking/returning `[]*typev1.Tag`). Pure map ops — clone, key removal, map+map merge — stay in stdlib `maps`/`slices` at the call site (`maps.Clone`, `maps.Copy`); do not add `Clone`/`Without` here.
- Dependencies: [gomatic/go-error](https://github.com/gomatic/go-error) for the sentinel mechanism (`ErrInvalidPair` is an `errs.Const`, matched with `errors.Is`; never `errors.New`, `fmt.Errorf` only to `%w`-wrap) and the **private** `github.com/skykernel/api` for the proto `Tag` type — so this repo must stay private until `skykernel/api` is public, even if the rest of `gomatic/go-*` goes public.
- Value-oriented, immutable, private by default. Gate: gofumpt, vet, staticcheck, govulncheck, gocognit ≤ 7, 100% coverage.
- `Makefile`, `.golangci.yaml`, `.editorconfig`, `.gitignore`, `.github/` are owned and pushed by `nicerobot/tools.repository` — do not edit in-tree; per-repo divergence goes in a `Makefile.local`.
