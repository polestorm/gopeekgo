//go:build nohighlight

package preview

import (
	"strings"
	"testing"
)

func TestPageWithoutHighlighting(t *testing.T) {
	out := string(Page([]byte("```go\nfunc main() {}\n```\n"), "x.md"))
	if strings.Contains(out, "hl-") || !strings.Contains(out, `<code class="language-go">func main() {}`) {
		t.Errorf("with the nohighlight tag a code block should stay plain text:\n%s", out)
	}
}
