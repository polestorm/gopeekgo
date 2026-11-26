# 0001 Scope and build

Status: accepted · 2026-11-25 (owner's request: a separate project depending on gomdgo and gomermaidgo, named
gopeekgo, Apache 2.0 licence)

## Context

The owner wants to look at any Markdown file as gomdgo and gomermaidgo render it, directly in Finder (Space).
Until now another extension (QLMarkdown) with another parser did this. macOS accepts a Quick Look extension only
as part of an application, written natively and running in the sandbox. The prototype was made in
`gomdgo/tools/quicklook` and moved here the same day, without a commit there.

Ways considered to connect the extension with the Go code:

1. **Running the `gomdgo-mermaid` program from the extension.** Rejected: the sandbox makes starting programs
   hard, and the family avoids it (gomdgo, ADR 0014).
2. **A C archive built from Go and linked into the extension** (`go build -buildmode=c-archive`). Chosen.

## Decision

1. A separate repository and module `github.com/polestorm/gopeekgo`, depending on released versions of the
   libraries, without `replace`. The project is a consumer of the public API, not an internal tool of gomdgo.
2. `preview.Page(src, title)` builds the page in pure Go and returns no error. gomdgo configuration: GFM, alerts,
   footnotes, heading identifiers; raw HTML written as text (the library's default policy).
3. The extension hands the system finished HTML (a data-based preview, `QLIsDataBasedPreview`), with no view of
   its own and no network access.
4. Built by a script, `build.sh` (now `macos/build.sh`, ADR 0002), without an Xcode project; signed locally ad
   hoc; installed into `~/Applications`.
5. Syntax highlighting will come from the gohighlightgo library. Originally an option added with the
   `highlight` tag; changed in ADR 0003 to on by default (the `nohighlight` tag turns it off).

## Consequences

- Images with relative paths are not loaded; to be done as attachments of the preview reply.
- Built for the architecture of the machine that runs the build script; distribution needs a Developer ID
  signature, notarization and a universal binary (a separate decision).
- The system picks one extension per file type; another Markdown extension has to be turned off.
- No CI at first: the tests of `preview` needed access to the private modules of the family (changed in ADR
  0002, once the libraries were public).
