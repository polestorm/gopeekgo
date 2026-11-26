package preview

import (
	"strings"
	"testing"
)

func TestPageRendersMarkdownAndDiagram(t *testing.T) {
	src := "# Title\n\n| a | b |\n|---|---|\n| 1 | 2 |\n\n```mermaid\nflowchart LR\n  A --> B\n```\n"
	out := string(Page([]byte(src), "file <1>.md"))
	for _, want := range []string{
		"<!DOCTYPE html>",
		"<title>file &lt;1&gt;.md</title>",
		"<h1 id=",
		"<table>",
		"<figure class=\"diagram\">",
		"<svg",
		"<footer>gopeekgo " + Version + " · " + libraries + " · ",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("%q missing from the output", want)
		}
	}
}

func TestPageKeepsGoingAfterDiagramError(t *testing.T) {
	src := "before\n\n```mermaid\nthis is not a <diagram>\n```\n\nafter\n"
	out := string(Page([]byte(src), "x.md"))
	for _, want := range []string{"<p>before</p>", "<p>after</p>", "class=\"diagram-error\"", "mermaid, line 3:", "&lt;diagram&gt;"} {
		if !strings.Contains(out, want) {
			t.Errorf("%q missing from the output:\n%s", want, out)
		}
	}
}

func TestPageEscapesRawHTML(t *testing.T) {
	out := string(Page([]byte("<script>alert(1)</script>\n"), "x.md"))
	if strings.Contains(out, "<script>") {
		t.Errorf("raw HTML reached the preview:\n%s", out)
	}
}
