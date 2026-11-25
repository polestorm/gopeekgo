# gopeekgo

Quick Look preview of Markdown files on macOS, rendered in pure Go: select a `.md` file in Finder, press
Space.

The page is produced by [gomdgo](https://github.com/polestorm/gomdgo) (Markdown → HTML, CommonMark + GFM)
and [gomermaidgo](https://github.com/polestorm/gomermaidgo) (` ```mermaid ` blocks → SVG). Both run inside
the Quick Look extension as a linked Go archive: no JavaScript, no helper process, no network.

Status: **working prototype** (2026-11-25). Local builds only, signed ad hoc.

## Install

Requires Go with cgo and the Xcode command line tools.

```sh
./build.sh install     # build, copy to ~/Applications/Gopeekgo.app, register the extension
./build.sh uninstall   # unregister and remove
go test ./...          # tests of the Go part (run anywhere)
```

macOS uses one preview extension per file type. If another Markdown preview is installed, turn it off in
System Settings (General → Login Items & Extensions → Quick Look), or for example:

```sh
pluginkit -e ignore -i org.sbarex.QLMarkdown.QLExtension   # off
pluginkit -e use    -i org.sbarex.QLMarkdown.QLExtension   # back on
qlmanage -r
```

## What the preview shows

- CommonMark and GFM (tables, task lists, strikethrough, autolinks), alerts (`> [!NOTE]`), footnotes,
  heading ids;
- Mermaid diagrams drawn as SVG at conversion time; a diagram with an error is replaced by the message,
  the line number and the source of the block;
- raw HTML from the document is shown as text, never executed;
- light and dark appearance;
- a footer with the size of the source, the size of the result and the conversion time.

Not yet: syntax highlighting of code blocks (planned through
[gohighlightgo](https://github.com/polestorm/gohighlightgo), as a build option), images with relative
paths, a universal binary.

## Comparison

First measurements, 2026-11-25, Apple silicon, against QLMarkdown 1.5.7. The advantage is
`(other − gopeekgo) / other × 100%`; a negative value means gopeekgo is worse. Method and caveats:
`docs/testing.md`.

| | gopeekgo | QLMarkdown | advantage |
|---|--:|--:|--:|
| application on disk | 6.6 MB | 64.0 MB | 90% |
| extension executable, arm64 | 6.68 MB | 5.43 MB | −23% |
| conversion, 24 KB document | 4.2 ms | 21.5 ms | 80% |
| conversion, 1.3 MB document | 14.9 ms | 199.8 ms | 93% |
| peak memory, 1.3 MB document | 31.8 MB | 60.3 MB | 47% |

QLMarkdown does more: it highlights code and has many settings. A fuller comparison with other previews
on the market is planned (`docs/roadmap.md`).

## Layout

| Path | Content |
|---|---|
| `preview/` | the HTML page of the preview; plain Go, no cgo |
| `bridge/` | two C functions over `preview`, built with `-buildmode=c-archive` |
| `macos/` | the Quick Look extension and the host application in Swift, property lists, entitlements |
| `build.sh` | builds, signs and installs the bundle without an Xcode project |

## Licence

Apache License 2.0 — see `LICENSE` and `NOTICE`.
