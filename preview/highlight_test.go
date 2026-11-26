//go:build !nohighlight

package preview

import (
	"strings"
	"testing"
)

func TestPageHighlightsCode(t *testing.T) {
	out := string(Page([]byte("```go\nfunc main() {}\n```\n\n```mermaid\nflowchart LR\n  A --> B\n```\n"), "x.md"))
	for _, want := range []string{`<pre class="hl"><code class="language-go">`, `<span class="hl-kw">func</span>`, ".hl-kw", "<svg"} {
		if !strings.Contains(out, want) {
			t.Errorf("brak %q w wyniku", want)
		}
	}
}

func TestPageKeepsUnknownLanguagePlain(t *testing.T) {
	out := string(Page([]byte("```no-such-language\n<b>x</b>\n```\n"), "x.md"))
	if !strings.Contains(out, "&lt;b&gt;x&lt;/b&gt;") || strings.Contains(out, "<b>x</b>") {
		t.Errorf("blok w nieznanym języku nie jest zwykłym, bezpiecznym tekstem:\n%s", out)
	}
}
