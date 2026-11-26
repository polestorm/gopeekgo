# Usage

[← README](../README.md)

## The Quick Look extension

Requires macOS 13 or newer, Go 1.22 or newer with cgo, and the Xcode command line tools. No Xcode project is
needed.

```sh
macos/build.sh install                     # build, copy to ~/Applications/GoPeekGo.app, register the extension
TAGS=nohighlight macos/build.sh install    # the same without syntax highlighting (about 0.5 MB smaller)
macos/build.sh uninstall                   # unregister and remove
qlmanage -p README.md                      # show a preview without Finder
```

The application, GoPeekGo, is built for the architecture of the machine and signed ad hoc, for use on that machine. Its
only job is to carry the extension; it shows a short note when opened. The footer of every preview ("GoPeekGo 0.1.0 · …") names the
version and the libraries, so you can tell which preview you are looking at.

macOS uses one preview extension per file type. If another Markdown preview is installed, turn it off in System
Settings (General → Login Items & Extensions → Quick Look), or for example:

```sh
pluginkit -m -A -p com.apple.quicklook.preview                # + means on, - means off
pluginkit -e ignore -i org.sbarex.QLMarkdown.QLExtension      # QLMarkdown off
pluginkit -e use    -i org.sbarex.QLMarkdown.QLExtension      # back on
qlmanage -r
```

## The command line

`gopeekgo-render` writes the page the extension returns, so the preview can be opened in a browser, compared or
measured. It runs anywhere Go runs; it needs no cgo.

```sh
go install github.com/polestorm/gopeekgo/cmd/gopeekgo-render@latest
gopeekgo-render notes.md > notes.html
gopeekgo-render < notes.md > notes.html
gopeekgo-render -version
go install -tags nohighlight github.com/polestorm/gopeekgo/cmd/gopeekgo-render@latest   # without highlighting
```

## What the page contains

- CommonMark 0.31.2 with the GFM extensions (tables, task lists, strikethrough, autolinks), GitHub alerts,
  footnotes and heading identifiers, from gomdgo.
- Mermaid diagrams drawn as SVG by gomermaidgo; a diagram with an error is replaced by the message, its line
  number and the source of the block.
- Code blocks in any of the 298 languages of gohighlightgo, coloured with its `auto` theme; a block in an
  unknown language stays plain text.
- Light and dark appearance, following the system.
- Raw HTML from the document is shown as text, and `javascript:` links are removed. The page has no scripts;
  its stylesheet and diagrams are inline.

Not yet: images with relative paths, math, YAML front matter as metadata, a universal binary.

## Go package

`github.com/polestorm/gopeekgo/preview` builds the page: `preview.Page(src, title)` returns a complete HTML
document and never fails. The API may change before version 1.0.
