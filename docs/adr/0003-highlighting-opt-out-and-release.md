# 0003 Syntax highlighting on by default; the release waits for gohighlightgo

Status: accepted · 2026-11-25 (owner's decision; changes point 5 of ADR 0001)

## Context

ADR 0001 planned highlighting as an option added at build time (the `highlight` tag). A preview is read by a
person, its HTML goes nowhere else, and the competitor (QLMarkdown) highlights by default. The gohighlightgo
library is being built in a separate repository.

## Decision

1. Syntax highlighting is **on by default** in gopeekgo (opt-out). A plain build includes gohighlightgo; the
   `nohighlight` tag turns it off: `preview/highlight.go` has `//go:build !nohighlight`, and
   `preview/nohighlight.go` (`//go:build nohighlight`) leaves the code blocks to gomdgo.
2. The build script gets a switch that passes `-tags nohighlight`. It also serves measurement: the size of the
   extension and the conversion time are measured in both variants (criterion H5 of gohighlightgo).
3. **The gopeekgo release waits until gohighlightgo is mature** (its acceptance criteria met and a version
   released). Until then gopeekgo stays a locally built prototype without a version tag.
4. `go.mod` points only to released versions of the libraries (ADR 0001, point 1), so the dependency on
   gohighlightgo is added only after its first release. Earlier trials are made locally with `go.work`, without
   a commit.

## Consequences

- The other consumers (`gomdgo/cmd/gomdgo-mermaid`, `gopressgo`) keep highlighting as an option; this decision
  concerns gopeekgo only. The gomdgo library itself does not know gohighlightgo.
- ADR 0001 of gohighlightgo describes the `highlight` tag for consumers; for gopeekgo this document applies.

## Update 2026-11-26

gohighlightgo v0.1.2 is released and public, so point 4 is met: the dependency is in `go.mod` and the adapter
from `gohighlightgo/docs/adapters/gopeekgo.patch` is applied (ADR 0004). The switch of point 2 is
`TAGS=nohighlight macos/build.sh`. Point 3 still holds: the version tag is the owner's decision.

## Update 2026-11-26: release v0.1.0 before the 24-hour fuzzing of gohighlightgo

Owner's decision: gopeekgo v0.1.0 is released and the repository made public without waiting for the 24-hour
fuzzing of gohighlightgo (`FuzzHTML`, `FuzzDefinition`), which started on 2026-11-26. gohighlightgo itself
released v0.1.0–v0.1.2 before that run (its ADR 0007). The result of the run will be added to
`docs/testing.md`; a fix it calls for comes as a new gohighlightgo version and a gopeekgo release that depends
on it. Point 3 is replaced by this update.
