# 0004 Getting ready for publication: version, highlighting, comparison, licences

Status: accepted · 2026-11-26 (owner's request: make the repository consistent with the public gomdgo,
gomermaidgo, gopdfgo and gohighlightgo before publishing it)

## Context

The owner wants to make gopeekgo public soon. The public repositories of the family have: a version number in the
comparison table and in their commands, English doc comments (pkg.go.dev shows them; gohighlightgo v0.1.2), a
repeatable comparison on frozen input with a described method, few files in the root, the Apache 2.0 licence and
a `NOTICE` that lists the parts by other authors. gohighlightgo has released v0.1.2 and is public, so the
condition of ADR 0003, point 4, for adding the dependency is met.

## Decision

1. **Version 0.1.0** in one place: `preview.Version`. The preview footer, `gopeekgo-render -version` and the
   heading of the README table show it; `CFBundleShortVersionString` in both `macos/*-Info.plist` files must be
   the same (`TestVersionMatchesBundles`). No `v0.1.0` tag without the owner's decision (ADR 0003, point 3).
2. **Highlighting** from gohighlightgo v0.1.2 in `go.mod`, on by default; the adapter from
   `gohighlightgo/docs/adapters/gopeekgo.patch`. `macos/build.sh` accepts `TAGS=nohighlight`.
3. **English doc comments** in the packages and commands pkg.go.dev shows (`preview`, `bridge`,
   `cmd/gopeekgo-render`, `tools/bench`, `tools/features`). (Extended to all comments and documents by ADR 0005.)
4. **Comparison with its input and method**: time, memory and size from `tools/bench` on `testdata/bench/` (with
   a new input full of code, `code.md`, since we highlight now), the feature rows from `tools/features` on
   `testdata/features/` (one file per feature, the recognition rule in the code). Both tools record QLMarkdown's
   settings, because the result depends on them. The method is in `docs/benchmark.md`, the results in
   `docs/testing.md`.
5. **Medals** only where one value is better whatever the reader prefers (time, memory, size, having a feature,
   no scripts, popularity, the way of distribution); ties and rows that are a matter of taste (how diagrams are
   drawn, raw HTML) get none. The advantage column in percent only on numbers.
6. **Licences**: `NOTICE` lists the family libraries and the Arimo and Noto Emoji typefaces (SIL OFL 1.1) that
   gomermaidgo embeds in our executable; `macos/build.sh` copies `LICENSE`, `NOTICE` and the OFL texts into
   `Gopeekgo.app/Contents/Resources`.
7. **A short README**, as in gohighlightgo: no "Status" section; installation and usage in `docs/usage.md`.

## Consequences

- The root keeps `LICENSE`, `NOTICE`, `README.md`, `go.mod`, `go.sum` and `.gitignore`. `go.sum` is needed
  because there are dependencies (as in gopdfgo); fewer files are not possible without giving up modules.
- Highlighting costs 497,696 B of the executable and 0.7–1.1 ms on a small file (H5, `docs/testing.md`); the
  executable size against QLMarkdown went from −23% to −32%.
- The texts of the host application were still in Polish; translated on 2026-11-26 (ADR 0005).
- `go install …@latest` from the README works only once the repository is public and tagged.
