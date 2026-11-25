# gopeekgo

Quick Look preview of Markdown files on macOS, rendered in pure Go: select a `.md` file in Finder and press
Space. The page is built by [gomdgo](https://github.com/polestorm/gomdgo) (CommonMark + GFM → HTML) and
[gomermaidgo](https://github.com/polestorm/gomermaidgo) (` ```mermaid ` blocks → SVG), linked into the
extension as a Go archive: no JavaScript, no helper process, no network.

```sh
macos/build.sh install   # build, copy to ~/Applications/Gopeekgo.app, register the extension
qlmanage -p README.md    # or Space in Finder
```

## Why gopeekgo

Compared with [QLMarkdown](https://github.com/sbarex/QLMarkdown), a popular free and open-source Markdown preview for
macOS (details in [Comparison](#comparison)):

- **Diagrams without a browser engine**: Mermaid blocks are drawn as SVG while the page is built; no script
  runs in the preview and nothing is fetched from the network.
- **Faster**: 80–93% less conversion time on the three test documents.
- **Smaller**: a 6.6 MB application instead of 64 MB, and half the memory on a 1.3 MB document.
- **Safe by default**: raw HTML from the document is shown as text and `javascript:` links are removed.
- **Shows its own errors**: a diagram that cannot be drawn is replaced by the message, its line number and the
  source of the block; the rest of the page is still shown.

Where it is weaker: no syntax highlighting yet (planned), no math, emoji shortcodes, front matter, relative
images or settings. It is built from source for the architecture of the machine and signed ad hoc, and its
extension executable is 23% larger, because it carries the Go runtime and everything else.

## Contents

- [Why gopeekgo](#why-gopeekgo) · [Comparison](#comparison) · [Status](#status) · [Install](#install) ·
  [Documentation](#documentation) · [How the comparison was measured](#how-the-comparison-was-measured) ·
  [License](#license)

## Comparison

QLMarkdown converts the document inside its extension and leaves diagrams and math to JavaScript libraries
(Mermaid, MathJax) running in the preview. It does more than gopeekgo and has many settings; the table shows
both the numbers and the features.

| | gopeekgo | QLMarkdown 1.5.7 | gopeekgo advantage[^adv] |
|---|:---:|:---:|:---:|
| Application on disk[^size] | 6.6 MB 🥇 | 64.0 MB | 90% |
| Extension executable, arm64 | 6.68 MB | 5.43 MB 🥇 | −23% |
| Conversion, 24 KB document[^bench] | 4.2 ms 🥇 | 21.5 ms | 80% |
| Conversion, 0.2 KB with a diagram | 2.9 ms 🥇 | 16.5 ms | 82% |
| Conversion, 1.3 MB document | 14.9 ms 🥇 | 199.8 ms | 93% |
| Peak memory, 0.2 KB with a diagram | 7.1 MB 🥇 | 18.2 MB | 61% |
| Peak memory, 1.3 MB document | 31.8 MB 🥇 | 60.3 MB | 47% |
| Mermaid diagrams[^probe] | SVG drawn during conversion 🥇 | Mermaid.js in the preview | |
| Diagram types | 24 of Mermaid's 33 | all of Mermaid 12.0 🥇 | |
| Network, JavaScript in the preview[^net] | none 🥇 | Mermaid and MathJax from a CDN, or embedded | |
| `<script>` in the source | written as text | removed | |
| Other raw HTML (`<b>`) | written as text | rendered | |
| `javascript:` link | removed 🥇 | kept in the page | |
| Tables, task lists, strikethrough, autolinks | yes | yes | |
| Footnotes, heading identifiers | yes | yes | |
| GitHub alerts (`> [!NOTE]`) | yes | yes, a setting | |
| Syntax highlighting | not yet | yes 🥇 | |
| Math | no | yes, MathJax 🥇 | |
| YAML front matter | no (read as Markdown) | a table 🥇 | |
| Emoji shortcodes (`:smile:`) | no | yes 🥇 | |
| Images with relative paths | no | yes, embedded 🥇 | |
| Settings | none | many, in its application 🥇 | |
| Light and dark appearance | yes | yes | |
| Architectures | the building machine's | universal (arm64 + x86_64) 🥇 | |
| Distribution | build from source, signed ad hoc | notarized, Developer ID 🥇 | |
| License | Apache 2.0 | GPL-3.0 | |
| GitHub popularity[^pop] | new | 3,795 stars, 124 forks | |

[^adv]: `(QLMarkdown − gopeekgo) / QLMarkdown × 100%`; a negative value means gopeekgo is worse. 🥇 marks the
    better value in a row.
[^size]: `du` of the installed application; the executable is the arm64 slice of QLMarkdown's universal file
    (11.4 MB in total). QLMarkdown keeps highlighting and Mermaid outside it (`libwrapper_highlight.dylib`
    6.9 MB, `mermaid.min.js.gz` 1.6 MB); gopeekgo has everything in its executable.
[^bench]: Whole process of each project's command-line converter, median of 15 runs, memory from
    `/usr/bin/time -l`; measured 2026-11-25 on one Apple-silicon laptop. gopeekgo's time includes drawing the
    diagram; QLMarkdown's does not, because Mermaid draws it later in the preview.
[^probe]: The feature rows come from converting one small probe file with both programs on 2026-11-25:
    QLMarkdown through `qlmarkdown_cli` with the settings of the installed application, gopeekgo through
    `preview.Page`. "No" means the feature did not appear in the output.
[^net]: With the settings used, QLMarkdown loads Mermaid 12.0.0 and MathJax 4.1.3 from `cdn.jsdelivr.net`;
    it can also embed them.
[^pop]: Stars and forks on 2026-11-25. The last release of QLMarkdown is from 2026-10-02.

## Status

Working prototype (2026-11-25). It depends on the released versions of the libraries
(gomdgo v0.1.1, gomermaidgo v0.1.0), so it sees them as any other consumer does; a rendering error found here
is fixed in the library, not worked around. Requires macOS 13 or newer, Go 1.22 or newer with cgo, and the
Xcode command line tools.

## Install

```sh
macos/build.sh install     # build, copy to ~/Applications/Gopeekgo.app, register the extension
macos/build.sh uninstall   # unregister and remove
go test ./...              # tests of the Go part (run anywhere)
```

macOS uses one preview extension per file type. If another Markdown preview is installed, turn it off in
System Settings (General → Login Items & Extensions → Quick Look), or for example:

```sh
pluginkit -e ignore -i org.sbarex.QLMarkdown.QLExtension   # off
pluginkit -e use    -i org.sbarex.QLMarkdown.QLExtension   # back on
qlmanage -r
```

## Documentation

In Polish: [`docs/testing.md`](docs/testing.md) (tests, measurements, comparison with QLMarkdown),
[`docs/adr/`](docs/adr/) (decisions), [`docs/roadmap.md`](docs/roadmap.md).

## How the comparison was measured

Both projects were run only as programs, on the same files and the same machine; the code of other previews and
Markdown parsers is not read. The conversion times come from the command-line converters
(`gomdgo-mermaid -standalone -alerts -footnotes -heading-ids` and `qlmarkdown_cli` from the QLMarkdown bundle),
not from the extensions, so they do not include the time from pressing Space to seeing the page; that is not
measured yet. The 24 KB document is `docs/testing.md` of gomdgo, the 1.3 MB document is 40 copies of it with
gomdgo's `README.md`. gopeekgo's numbers were measured on the prototype before it moved to this repository,
with the gomdgo code of that day; the method and caveats are in [`docs/testing.md`](docs/testing.md).

## License

Apache License 2.0, see `LICENSE` and `NOTICE`.
