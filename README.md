# gopeekgo

Quick Look preview of Markdown files on macOS, rendered in pure Go: select a `.md` file in Finder and press
Space. [gomdgo](https://github.com/polestorm/gomdgo) turns CommonMark and GFM into HTML,
[gomermaidgo](https://github.com/polestorm/gomermaidgo) draws ` ```mermaid ` blocks as SVG and
[gohighlightgo](https://github.com/polestorm/gohighlightgo) colours the code, all inside the extension: no
JavaScript, no helper process, no network.

```sh
macos/build.sh install                                                # the Quick Look extension
go install github.com/polestorm/gopeekgo/cmd/gopeekgo-render@latest   # the same page from the command line
```

## Why gopeekgo

Compared with [QLMarkdown](https://github.com/sbarex/QLMarkdown), a popular free and open-source Markdown
preview for macOS (details in [Comparison](#comparison)):

- **Fast**: 80–99% less time than QLMarkdown's converter on the five test inputs; 31 ms instead of 2.2 s for
  1.3 MB of code.
- **Small**: a 7.3 MB application instead of 64 MB, and less memory on every input.
- **No scripts in the preview**: diagrams are drawn as SVG and code is highlighted while the page is built, so
  the page has no `<script>` and loads no library from the network.
- **Safe by default**: raw HTML from the document is shown as text and `javascript:` links are removed.
- **Shows its own errors**: a diagram that cannot be drawn is replaced by the message, its line number and the
  source of the block; the rest of the page is still shown.

Where it is weaker: no math, emoji shortcodes, YAML front matter, images with relative paths or settings, and
fewer Mermaid diagram types. It is built from source for the architecture of the machine and signed ad hoc, and
its extension executable is 32% larger, because it carries the Go runtime, the fonts for measuring diagram text
and everything else.

## Contents

- [Why gopeekgo](#why-gopeekgo) · [Comparison](#comparison) · [Documentation](#documentation) ·
  [License](#license)
- [Usage](docs/usage.md): installing, turning another Markdown preview off, the command line, building without
  highlighting
- [How the comparison was measured](docs/benchmark.md): the input files, the two commands behind the table and
  what the numbers do not show

## Comparison

| | gopeekgo 0.1.0 | QLMarkdown 1.5.7 | gopeekgo advantage[^adv] |
|---|:---:|:---:|:---:|
| Application on disk[^size] | 7.3 MB 🥇 | 63.9 MB | 89% |
| Extension executable, arm64 | 7.2 MB | 5.4 MB 🥇 | −32% |
| Time, 24 kB document[^bench] | 3.3 ms 🥇 | 21.4 ms | 84% |
| Time, 0.2 kB with a diagram | 3.3 ms 🥇 | 16.7 ms | 80% |
| Time, 1.3 MB document | 15.7 ms 🥇 | 213.6 ms | 93% |
| Time, 15 kB of code in 19 languages | 3.5 ms 🥇 | 47.9 ms | 93% |
| Time, 1.3 MB of code | 31 ms 🥇 | 2,222 ms | 99% |
| Peak memory, 24 kB document | 8.1 MB 🥇 | 21.7 MB | 63% |
| Peak memory, 1.3 MB document | 32.8 MB 🥇 | 59.5 MB | 45% |
| Peak memory, 1.3 MB of code | 42.8 MB 🥇 | 57.2 MB | 25% |
| Tables, task lists, strikethrough, `www.` links, footnotes, heading identifiers[^feat] | yes | yes | |
| Syntax highlighting | yes | yes | |
| GitHub alerts (`> [!NOTE]`) | yes | a setting, off in the settings measured | |
| Mermaid diagrams | drawn as SVG while the page is built | drawn by Mermaid.js in the preview | |
| Mermaid diagram types[^types] | 24 of 33 | all of Mermaid 12.0 🥇 | |
| Scripts in the page | none 🥇 | Mermaid.js and MathJax, from `cdn.jsdelivr.net` or embedded (a setting) | |
| Math | no | yes 🥇 | |
| YAML front matter | read as Markdown | shown as metadata 🥇 | |
| Emoji shortcodes (`:smile:`) | no | yes 🥇 | |
| Image with a relative path[^img] | not shown | embedded in the page 🥇 | |
| `<script>` in the source | shown as text | shown as text | |
| Other raw HTML (`<b>`) | shown as text | rendered | |
| `javascript:` link | removed 🥇 | kept | |
| Settings | none | many, in its application 🥇 | |
| Architectures | the building machine's | universal (arm64 + x86_64) 🥇 | |
| Distribution | built from source, signed ad hoc | notarized, Developer ID 🥇 | |
| License | Apache 2.0 | GPL 3.0 | |
| GitHub popularity[^pop] | new | 3,795 stars, 124 forks 🥇 | |

To repeat the comparison, see [docs/benchmark.md](docs/benchmark.md). If you maintain QLMarkdown and get
different results, please open an issue.

[^adv]: `(QLMarkdown − gopeekgo) / QLMarkdown × 100%`; a negative value means gopeekgo is worse. 🥇 marks the
    better value where one is better whatever the reader prefers; rows where the two are equal or where it is a
    matter of taste have none.
[^size]: Sum of the file sizes in the installed application. The executable row compares the arm64 slice of
    QLMarkdown's universal file (11.4 MB in total); QLMarkdown keeps highlighting and Mermaid outside it
    (`libwrapper_highlight.dylib` 6.9 MB, `mermaid.min.js.gz` 1.6 MB), gopeekgo has everything in its executable.
[^bench]: Whole process of each project's command-line converter (`gopeekgo-render`, which writes exactly the
    page the extension returns, and `qlmarkdown_cli` from the QLMarkdown bundle) on the frozen files in
    [`testdata/bench/`](testdata/bench/); median of 15 runs, peak memory from `getrusage`. `go run ./tools/bench`,
    2026-11-26, Apple M4. gopeekgo's time includes drawing the diagram; QLMarkdown's does not, because
    Mermaid.js draws it later in the preview. [More](docs/benchmark.md)
[^feat]: The feature rows come from `go run ./tools/features`: one small file per feature in
    [`testdata/features/`](testdata/features/), converted by both programs, and a fixed rule that reads the page.
    QLMarkdown runs with the settings of its installed application, recorded with the results in
    [`docs/testing.md`](docs/testing.md).
[^types]: Mermaid 12.2 recognises 33 diagram types; gomermaidgo draws 24 of them
    ([its comparison](https://github.com/polestorm/gomermaidgo#comparison)). QLMarkdown runs Mermaid 12.0 itself.
[^img]: gopeekgo leaves the address in the page, and a Quick Look preview built from data cannot load it.
[^pop]: Stars and forks on 2026-11-25. The last release of QLMarkdown is from 2026-10-02. Architectures,
    distribution and license were checked with `lipo`, `spctl` and the GitHub API on the same day.

## Documentation

- [`docs/usage.md`](docs/usage.md): installing and using the preview and the command.
- [`docs/benchmark.md`](docs/benchmark.md): how the comparison is measured, with the input files.
- [`docs/testing.md`](docs/testing.md): tests and the results of every measurement.
- [`docs/adr/`](docs/adr/): decisions; [`docs/roadmap.md`](docs/roadmap.md): what is done and what is next.

## License

Apache License 2.0, see `LICENSE` and `NOTICE`. The built application and command include the Arimo and Noto
Emoji typefaces through gomermaidgo, under the SIL Open Font License 1.1; `NOTICE` says where their license
texts are, and `macos/build.sh` puts them into the application.
