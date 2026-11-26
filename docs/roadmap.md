# Roadmap

[← README](../README.md)

| # | Scope | State |
|---|---|---|
| 0 | prototype: preview page, C bridge, extension, `macos/build.sh`, local install | done (2026-11-25) |
| 1 | look at real documents: dark mode, long tables, wide diagrams, large files; list of findings for the libraries | not started |
| 2 | images with relative paths (attachments of the preview reply) | not started |
| 3 | a switch in `macos/build.sh` to build with local copies of the libraries (`go.work`) | not started |
| 4 | syntax highlighting from gohighlightgo, on by default, the `nohighlight` tag turns it off (ADR 0003); size and time measured in both variants; a condition of the release | done with gohighlightgo v0.1.2, H5 measured (2026-11-26); released in v0.1.0 before gohighlightgo's 24-hour fuzzing (ADR 0003, update); its result to be added |
| 5 | comparison with other previews on the market (QLMarkdown, the system preview, others): size, time, memory, features; table in `README.md` | QLMarkdown, repeatable (`tools/bench`, `tools/features`, 2026-11-26); other previews to be added |
| 6 | time from pressing Space to seeing the page (not only conversion) | not started |
| 7 | distribution: universal binary, Developer ID, notarization | to be decided |
| 8 | CI: `go vet`, `go test`, building the application on macOS (ADR 0002) | disabled on 2026-11-25 (billing of the private repository), enabled again on 2026-11-26 when the repository became public |
| 9 | publication: version 0.1.0, English doc comments, `NOTICE` with the OFL fonts, comparison with its input (ADR 0004) | done: public and tagged `v0.1.0` (2026-11-26) |
| 10 | texts of the host application in English for publication | done (2026-11-26) |
| 11 | documentation in English: documents and comments (ADR 0005) | done (2026-11-26) |
