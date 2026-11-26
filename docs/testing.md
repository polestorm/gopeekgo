# Tests and measurements

[← README](../README.md)

## Tests

- `go test ./...`: package `preview` (the page contains the Markdown output and the diagram; a diagram error does
  not stop the preview; raw HTML does not reach the page; code blocks are highlighted and a block in an unknown
  language stays safe plain text; `preview.Version` equals `CFBundleShortVersionString` in both `Info.plist`
  files).
- `go test ./tools/bench`: the input of the comparison is byte for byte the one listed in
  `testdata/SHA256SUMS`, and every file in `testdata/bench` and `testdata/features` is listed.
- `go test -tags nohighlight ./...`: the same without highlighting (the test checks that a code block is plain
  text).
- The macOS part is checked by hand: `macos/build.sh install`, then `qlmanage -p file.md` or Space in Finder. The
  footer of the page ("GoPeekGo 0.1.0 · gomdgo + gomermaidgo + gohighlightgo · size · time") shows that the
  preview comes from this extension, and from which variant.

## Comparison measurements: gopeekgo and QLMarkdown

The measurement can be repeated: `go run ./tools/bench` (time, memory, size) and `go run ./tools/features` (the
feature rows) in the root of the repository, on the files in `testdata/bench/` and `testdata/features/`. The
input (with checksums), the method and the limits are described in [`benchmark.md`](benchmark.md). Every number
of the comparison in `README.md` must come from a row of this section; measurements made in another way are not
published. The advantage is `(QLMarkdown − gopeekgo) / QLMarkdown × 100%`.

### 2026-11-26, version 0.1.0, QLMarkdown 1.5.7 with fixed flags, application built by the tool

The first run in which the results do not depend on the machine's setup: QLMarkdown gets the fixed flags of
`tools/internal/rival` (its saved settings plus `--alert on`, appearance light), the size rows measure the
application `tools/bench` builds from the tree, and `testdata/SHA256SUMS` pins the input. Environment: Apple M4,
macOS 27.0.1, mains power, Go 1.26.6, gomdgo v0.1.1, gomermaidgo v0.1.0, gohighlightgo v0.1.2, 15 runs per cell.
This section is the source of the README table.

| input | size | gopeekgo, median (min) | QLMarkdown, median (min) | advantage | gopeekgo memory | QLMarkdown memory | advantage |
|---|--:|--:|--:|--:|--:|--:|--:|
| `document.md` | 24,499 B | 3.7 ms (3.3 ms) | 21.4 ms (21.0 ms) | 82% | 8.0 MB | 21.8 MB | 63% |
| `diagram.md` | 202 B | 3.3 ms (3.3 ms) | 17.0 ms (16.3 ms) | 80% | 8.8 MB | 18.3 MB | 52% |
| 55 × `document.md` | 1,347,445 B | 15.9 ms (15.3 ms) | 212.1 ms (206.7 ms) | 93% | 32.8 MB | 59.6 MB | 45% |
| `code.md` | 15,429 B | 3.7 ms (3.4 ms) | 48.3 ms (47.5 ms) | 92% | 8.6 MB | 21.2 MB | 59% |
| 87 × `code.md` | 1,342,323 B | 30.5 ms (30.2 ms) | 2,207.3 ms (2,195.1 ms) | 99% | 42.9 MB | 57.2 MB | 25% |

| size on disk | gopeekgo | QLMarkdown | advantage |
|---|--:|--:|--:|
| application, sum of file sizes | 7,259,274 B | 63,896,352 B | 89% |
| extension bundle `.appex` | 7,177,500 B | 11,453,601 B | 37% |
| extension executable, arm64 slice | 7,174,048 B | 5,426,192 B | −32% |

The same with `-tags nohighlight` (criterion H5): 2.7, 2.6, 14.4, 2.4 and 10.2 ms; 6.3, 7.2, 30.5, 5.9 and
18.2 MB; application 6,761,578 B, executable 6,676,352 B. Highlighting therefore adds 497,696 B to the
executable, 0.7–1.3 ms on the small inputs and 20.3 ms on 1.3 MB of code, as in the run below.

Feature probe in the same run (`tools/features`): as in the table of the feature section, except GitHub alerts,
which QLMarkdown now renders (`yes`), because the fixed flags turn them on.

### 2026-11-26, version 0.1.0 with highlighting (gohighlightgo v0.1.2), working tree on `a6ff34b`

Environment: Apple M4, macOS 27.0.1, mains power, Go 1.26.6, gomdgo v0.1.1, gomermaidgo v0.1.0, gohighlightgo
v0.1.2, QLMarkdown 1.5.7 with the settings below, 15 runs per cell. New input: `code.md` (21 code blocks in 19
languages) and 87 copies of it; it shows the cost of highlighting, which the earlier files hardly had.

| input | size | gopeekgo, median (min) | QLMarkdown, median (min) | advantage | gopeekgo memory | QLMarkdown memory | advantage |
|---|--:|--:|--:|--:|--:|--:|--:|
| `document.md` | 24,499 B | 3.3 ms (3.2 ms) | 21.4 ms (20.8 ms) | 84% | 8.1 MB | 21.7 MB | 63% |
| `diagram.md` | 202 B | 3.3 ms (3.2 ms) | 16.7 ms (16.3 ms) | 80% | 8.8 MB | 18.2 MB | 52% |
| 55 × `document.md` | 1,347,445 B | 15.7 ms (15.2 ms) | 213.6 ms (210.5 ms) | 93% | 32.8 MB | 59.5 MB | 45% |
| `code.md` | 15,429 B | 3.5 ms (3.4 ms) | 47.9 ms (47.3 ms) | 93% | 8.7 MB | 21.1 MB | 59% |
| 87 × `code.md` | 1,342,323 B | 31.2 ms (30.6 ms) | 2,222.1 ms (2,201.0 ms) | 99% | 42.8 MB | 57.2 MB | 25% |

| size on disk | gopeekgo | QLMarkdown | advantage |
|---|--:|--:|--:|
| application, sum of file sizes | 7,258,656 B | 63,896,352 B | 89% |
| extension bundle `.appex` | 7,177,500 B | 11,453,601 B | 37% |
| extension executable, arm64 slice | 7,174,048 B | 5,426,192 B | −32% |

QLMarkdown's settings during the measurement (`qlmarkdown_cli --show-settings`; the tools record them too):
highlighting on, Mermaid 12.0.0 and MathJax 4.1.3 loaded from `cdn.jsdelivr.net`, GitHub alerts (`--alert`)
off, emoji as glyphs, YAML header for all files, local images embedded in the page, `--tag-filter` on.

#### gohighlightgo criterion H5: the same measurement with `-tags nohighlight`

| | with highlighting | `nohighlight` | difference |
|---|--:|--:|--:|
| extension executable, arm64 | 7,174,048 B | 6,676,352 B | +497,696 B (+7.5%) |
| application, sum of file sizes | 7,258,656 B | 6,761,578 B | +497,078 B |
| `document.md` | 3.3 ms | 2.6 ms | +0.7 ms |
| `code.md` | 3.5 ms | 2.4 ms | +1.1 ms |
| 87 × `code.md` (1.3 MB of code) | 31.2 ms | 10.5 ms | +20.7 ms |
| memory, 87 × `code.md` | 42.8 MB | 18.2 MB | +24.6 MB |

On small files the cost is mostly start-up (language definitions loaded lazily, names registered in the
`CodeBlocks` map). The `nohighlight` variant against QLMarkdown: 88%, 85%, 93%, 95%, 100% less time and 71%,
61%, 49%, 72%, 68% less memory, in the same order of inputs; executable −23%.

### 2026-11-26, commit `a6ff34b` with the measuring tool (without highlighting)

Environment: Apple M4, macOS 27.0.1, mains power, Go 1.26.6, gomdgo v0.1.1, gomermaidgo v0.1.0, QLMarkdown
1.5.7, 15 runs per cell.

| input | size | gopeekgo, median (min) | QLMarkdown, median (min) | advantage | gopeekgo memory | QLMarkdown memory | advantage |
|---|--:|--:|--:|--:|--:|--:|--:|
| `document.md` | 24,499 B | 3.6 ms (2.9 ms) | 22.4 ms (21.9 ms) | 84% | 6.3 MB | 21.7 MB | 71% |
| `diagram.md` | 202 B | 2.7 ms (2.6 ms) | 17.5 ms (17.1 ms) | 85% | 7.0 MB | 18.2 MB | 62% |
| 55 × `document.md` | 1,347,445 B | 14.4 ms (14.1 ms) | 214.3 ms (211.3 ms) | 93% | 30.6 MB | 59.4 MB | 49% |

| size on disk | gopeekgo | QLMarkdown | advantage |
|---|--:|--:|--:|
| application, sum of file sizes | 6,739,672 B | 63,896,352 B | 89% |
| extension bundle `.appex` | 6,679,628 B | 11,453,601 B | 42% |
| extension executable, arm64 slice | 6,676,176 B | 5,426,192 B | −23% |

### First measurement, 2026-11-25 (replaced)

The first numbers (4.2 ms / 21.5 ms on 24 kB; 14.9 ms / 199.8 ms on 1.3 MB; application 6.6 MB / 64.0 MB) came
from an ad hoc script, from the `gomdgo-mermaid` command instead of the code of this repository, and from files
outside the repository, so they could not be repeated. They were replaced by the measurement of 2026-11-26 and are
no longer quoted anywhere.

## Feature probe: gopeekgo and QLMarkdown

`go run ./tools/features`: each feature is a separate frozen file in `testdata/features/`, and the rule by which
the tool recognises the result in the page is written in `tools/features/main.go`. The rules read the content of
the page without its `<style>` and `<script>` elements, so that a class name in a stylesheet does not pass for
content (it did in the first version of the tool: QLMarkdown's alert came out "yes" because of its CSS).

### 2026-11-26, version 0.1.0, QLMarkdown 1.5.7 with its saved settings (before the fixed flags)

| feature | file | gopeekgo | QLMarkdown |
|---|---|---|---|
| tables, task lists, strikethrough, `www.` links, footnotes, heading identifiers | `table.md` … `heading-id.md` | yes | yes |
| GitHub alerts (`> [!NOTE]`) | `alert.md` | yes | no (with these settings; `--alert on` turns them on) |
| Mermaid diagrams | `mermaid.md` | SVG in the page | Mermaid.js in the preview |
| syntax highlighting | `code.md` | yes | yes |
| math | `math.md` | no | yes, a script in the preview |
| YAML front matter | `front-matter.md` | read as Markdown | shown as metadata |
| emoji shortcodes | `emoji.md` | no | yes |
| image with a relative path | `image.md` | left as an address | embedded in the page |
| `<script>` in the source | `raw-script.md` | shown as text | shown as text |
| other raw HTML (`<b>`) | `raw-tag.md` | shown as text | rendered |
| `javascript:` link | `javascript-link.md` | removed | kept |
| scripts in the pages, where from | all | none | 4 `<script>` elements; `cdn.jsdelivr.net` |

Outside the tool, checked by hand on 2026-11-25 (`codesign`, `spctl`, `lipo`, GitHub API): the QLMarkdown bundle
is signed with a Developer ID and notarized, the extension executable is universal (x86_64 + arm64), the licence
is GPL-3.0, 3,795 stars and 124 forks. Correction to the probe of 2026-11-25: QLMarkdown does not remove
`<script>`, it shows it as text (`&lt;script>`), like gopeekgo.

A finding for gomdgo (for the next batch of notes): a YAML header (`---`, `title: front`, `---`) comes out as a
thematic break and an `<h2>` heading. That follows CommonMark, but YAML headers are common in notes and
documentation; to be decided whether the library should skip them or return them separately (an option), or the
preview page should do it.
