# 0002 Repository layout, README and CI following the family libraries

Status: accepted · 2026-11-25 (owner's request: README and layout as in the public gomdgo, gomermaidgo and
gopdfgo)

## Context

The public repositories of the family share a pattern. The root holds only `LICENSE`, `NOTICE`, `README.md`,
`go.mod` (and `go.sum` when there are dependencies) and `.gitignore`; everything else is in directories. Each has
`.github/workflows/test.yml`. The English README has a fixed order: a few sentences on what the project does, a
short example, "Why …" with the advantages as bullets and a paragraph on weaknesses, contents, a comparison table
with one to three competitors (🥇 at the best value, footnotes with the method and the date), "Status",
"Documentation", "How the comparison was measured", "License".

ADR 0001 put CI off because the family modules were private. gomdgo and gomermaidgo are public now; on 2026-11-25
`go mod download` through `proxy.golang.org` and `sum.golang.org` gave checksums equal to `go.sum`, and
`go vet ./...` and `go test ./...` passed without private access.

## Decision

1. `build.sh` moves to `macos/build.sh` (it builds the macOS part; the script changes to the root of the
   repository, so the `build/` directory stays where it was).
2. README following the family pattern, with one competitor (QLMarkdown). Because the project's rule asks for the
   advantage in percent, the table has an extra column, "gopeekgo advantage", on the numeric rows.
3. CI as in the libraries: Go 1.22 and the latest, Linux and macOS; `gofmt`, `go vet`, `go test`, and on macOS
   building the application (`macos/build.sh` without installing).

## Consequences

- The point of ADR 0001 about no CI no longer applies. (The workflow was disabled on GitHub on 2026-11-25 until
  the next billing period; `docs/roadmap.md`, item 8.)
- The commands in the documentation read `macos/build.sh install` instead of `./build.sh install`.
- A comparison with other previews (`docs/roadmap.md`, item 5) adds columns to the same table.
- Later changes to the README layout: ADR 0004.
