# 0005 Documents and comments in English

Status: accepted · 2026-11-26 (owner's decision before publication)

## Context

Until now the README, `docs/benchmark.md`, `docs/usage.md` and the doc comments were in English, and the other
documents (`docs/testing.md`, `docs/roadmap.md`, the ADRs) and the remaining comments were in Polish, as in the
other libraries of the family. For a public repository, Polish text limits who can read and review the project,
and comments in the source are visible on pkg.go.dev and GitHub.

## Decision

1. All documents of the repository are in English: `README.md`, everything in `docs/`, including the ADRs.
2. All comments in the source are in English: Go (also unexported code and tests), Swift, the C header, the
   build script, `.gitignore`. Messages of the build script and of the tests are in English too.
3. Not covered:
   - the frozen benchmark input `testdata/bench/document.md` (a real document in Polish) and the Polish words in
     `testdata/bench/code.md` (part of gohighlightgo's corpus): changing them changes the benchmark;
   - the texts the host application shows to the user (a separate decision, `docs/roadmap.md`, item 10);
   - the instructions for Claude Code in `.claude/CLAUDE.md`, which are for the owner's tools, not for readers;
   - the copyright holder's name in `NOTICE`.

## Consequences

- This repository departs from the family rule "comments and documentation in Polish"; new documents and
  comments are written in English from now on.
