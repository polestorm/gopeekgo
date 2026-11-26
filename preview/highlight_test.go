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
			t.Errorf("%q missing from the output", want)
		}
	}
}

func TestPageKeepsUnknownLanguagePlain(t *testing.T) {
	out := string(Page([]byte("```no-such-language\n<b>x</b>\n```\n"), "x.md"))
	if !strings.Contains(out, "&lt;b&gt;x&lt;/b&gt;") || strings.Contains(out, "<b>x</b>") {
		t.Errorf("a block in an unknown language is not plain, safe text:\n%s", out)
	}
}
