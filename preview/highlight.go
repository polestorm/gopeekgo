//go:build !nohighlight

// Syntax highlighting of code blocks with gohighlightgo, on by default; the nohighlight build tag turns it
// off (nohighlight.go).
package preview

import (
	"html"

	"github.com/polestorm/gohighlightgo/highlight"
	gomdgo "github.com/polestorm/gomdgo/markdown"
)

// libraries names the libraries that build the page, for the footer.
const libraries = "gomdgo + gomermaidgo + gohighlightgo"

// highlightCSS is the "auto" theme: light or dark with the system appearance, like the rest of the page.
var highlightCSS, _ = highlight.ThemeCSS("auto")

// addHighlighting registers every language name gohighlightgo knows; "mermaid" stays with the diagrams.
func addHighlighting(blocks map[string]gomdgo.CodeBlockFunc) {
	for _, name := range highlight.Languages() {
		if blocks[name] == nil {
			blocks[name] = highlighted
		}
	}
}

// highlighted writes a block the way gomdgo does (<pre><code class="language-…">), with the class hl that
// the theme applies to.
func highlighted(b gomdgo.CodeBlock) ([]byte, error) {
	code, _ := highlight.HTML(b.Lang, b.Code)
	out := []byte(`<pre class="hl"><code class="language-` + html.EscapeString(b.Lang) + `">`)
	out = append(out, code...)
	return append(out, "</code></pre>\n"...), nil
}
