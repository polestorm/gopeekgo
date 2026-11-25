package preview

import (
	"strings"
	"testing"
)

func TestPageRendersMarkdownAndDiagram(t *testing.T) {
	src := "# Tytuł\n\n| a | b |\n|---|---|\n| 1 | 2 |\n\n```mermaid\nflowchart LR\n  A --> B\n```\n"
	out := string(Page([]byte(src), "plik <1>.md"))
	for _, want := range []string{
		"<!DOCTYPE html>",
		"<title>plik &lt;1&gt;.md</title>",
		"<h1 id=",
		"<table>",
		"<figure class=\"diagram\">",
		"<svg",
		"<footer>gopeekgo · gomdgo + gomermaidgo",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("brak %q w wyniku", want)
		}
	}
}

func TestPageKeepsGoingAfterDiagramError(t *testing.T) {
	src := "przed\n\n```mermaid\nto nie jest <diagram>\n```\n\npo\n"
	out := string(Page([]byte(src), "x.md"))
	for _, want := range []string{"<p>przed</p>", "<p>po</p>", "class=\"diagram-error\"", "mermaid, line 3:", "&lt;diagram&gt;"} {
		if !strings.Contains(out, want) {
			t.Errorf("brak %q w wyniku:\n%s", want, out)
		}
	}
}

func TestPageEscapesRawHTML(t *testing.T) {
	out := string(Page([]byte("<script>alert(1)</script>\n"), "x.md"))
	if strings.Contains(out, "<script>") {
		t.Errorf("surowy HTML przeszedł do podglądu:\n%s", out)
	}
}
