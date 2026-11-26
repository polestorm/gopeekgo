# How the comparison is measured

Every number and every feature row in the comparison table of the README can be repeated with one command on
the files in this repository. If you maintain one of the compared projects and get different results, please open an issue.

## Repeat it

Requirements: macOS, Go 1.22 or newer, this repository, and the other project installed (QLMarkdown puts its
command-line converter in `/Applications/QLMarkdown.app/Contents/Resources/qlmarkdown_cli`) and the Xcode
command line tools: for the size rows `tools/bench` builds the gopeekgo application itself with
`macos/build.sh` (same build tags, not installed, removed afterwards), so what is installed does not matter.

```sh
go run ./tools/bench                      # 15 runs per cell
go run ./tools/bench -runs 30             # more runs
go run ./tools/bench -rival /path/to/cli  # another location or another converter that takes a file name
go run ./tools/bench -tags nohighlight    # gopeekgo built without syntax highlighting
go run ./tools/features                   # the feature rows (a few seconds)
```

`bench` prints the environment (machine, system, power source, Go version, gopeekgo version and commit, build
tags, library versions, version and arguments of the other project) and two tables in Markdown; `features`
prints the environment and one table. Without the other project both report gopeekgo alone.

## Input

| File | Size | SHA-256 (first 16 digits) | What it is |
|---|--:|---|---|
| [`testdata/bench/document.md`](../testdata/bench/document.md) | 24,499 B | `516a6dd00a232218` | a real document: `docs/testing.md` of gomdgo at commit `587ee59`, prose with many tables, in Polish |
| [`testdata/bench/diagram.md`](../testdata/bench/diagram.md) | 202 B | `a591a71c80e40667` | a heading, a table, a task list, an alert and one Mermaid flowchart |
| [`testdata/bench/code.md`](../testdata/bench/code.md) | 15,429 B | `5171d8f8e55e7ced` | 21 fenced code blocks in 19 languages: the comparison corpus of [gohighlightgo](https://github.com/polestorm/gohighlightgo) v0.1.2 (`tools/compare/inputs/corpus`), without the first line of each file |
| large document | 1,347,445 B | generated | 55 copies of `document.md` joined, written to a temporary directory by the tool |
| large code document | 1,342,323 B | generated | 87 copies of `code.md` joined, the same way |
| [`testdata/features/`](../testdata/features/) | 17 files | | one small file per feature, read by `tools/features` |

The files are frozen: [`testdata/SHA256SUMS`](../testdata/SHA256SUMS) lists the SHA-256 of every file in
`testdata/bench` and `testdata/features`, and `go test ./tools/bench` fails when a file differs or is missing
from the list (`cd testdata && shasum -a 256 -c SHA256SUMS` checks the same by hand). A change to the input is
a change of the benchmark: it updates the list and gets a new section in `docs/testing.md`.

## Method

- **What runs.** For gopeekgo: `cmd/gopeekgo-render`, a command that writes exactly the page the Quick Look
  extension returns (`preview.Page`), built by the tool with `go build -trimpath`. For the other project: its
  own command-line converter, as a black box, with a fixed set of flags and the file name (below). Its source
  is not read.
- **Time** is the wall-clock time of the whole process, from start to exit, with the output collected through
  a pipe. Each cell is the median of the runs; the minimum is printed beside it. One extra run before the
  measured ones warms the file cache and is discarded.
- **Peak memory** is the maximum resident set size the system reports for the process (`getrusage`), median
  of the runs.
- **Size on disk** is the sum of the sizes of the regular files in the bundle: for gopeekgo the one the tool
  builds, for the other project the installed one. The executable row
  compares one architecture: for a universal file it takes the slice of the machine that runs the tool.
- **Features.** `tools/features` converts each file in `testdata/features` with both programs, from the file's
  directory (so that a relative image can be found), and applies a fixed rule to the page: for example "the
  page contains `<del>struck</del>`" for strikethrough. The rules are in `tools/features/main.go`. They read the
  page without its `<style>` and `<script>` elements, so a class name in a stylesheet does not count as
  content; only the math rule and the script count look at the scripts.
- **Settings of the other project.** QLMarkdown's converter would otherwise read the settings saved in its
  application, and they change what it does (highlighting, Mermaid, alerts, emoji). Both tools pass every
  setting that has a flag, the same on every machine: the list is
  [`tools/internal/rival/rival.go`](../tools/internal/rival/rival.go), and both tools print it with their
  results. It is the setup of a QLMarkdown 1.5.7 installation with GitHub alerts turned on, because gopeekgo
  renders them too.
- **Advantage** is `(other − gopeekgo) / other × 100%`. A negative value means gopeekgo is worse; it is
  published too.

## What the numbers do not show

- They measure conversion by the command-line programs, not the time from pressing Space to seeing the page.
  That is dominated by starting the extension process and the web view and is not measured yet.
- gopeekgo's time includes drawing the Mermaid diagram. QLMarkdown writes a `<script>` element and lets
  Mermaid.js draw the diagram later in the preview, so its real cost for `diagram.md` is higher than shown.
- Both highlight code, but differently (QLMarkdown through the C++ program highlight, gopeekgo through
  gohighlightgo), so the output differs; the times compare the work each one does on the same text. The cost of
  highlighting alone is the difference between `tools/bench` and `tools/bench -tags nohighlight`.
- A feature row says what appears in the page for one small file, not how complete the feature is.
- Process start is a large part of the two small inputs. The large document shows the cost of the conversion
  itself.
- One machine, one day. Results of each run, with the environment, are in [`testing.md`](testing.md).
