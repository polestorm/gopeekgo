// Package preview builds the HTML page of the gopeekgo Quick Look preview for a Markdown file. gomdgo
// converts the document to HTML, gomermaidgo draws the code blocks marked "mermaid" as SVG, and
// gohighlightgo colours the other code blocks (unless the program is built with the nohighlight tag).
//
// The package needs neither cgo nor macOS, so the page can be built and tested anywhere; the Quick Look
// extension calls it through the bridge command, and cmd/gopeekgo-render writes the same page from the
// command line.
//
// A preview is used to look at any file, including broken ones, so Page never fails: a diagram that cannot
// be drawn is replaced by its error message and source, and a panic in a library becomes a page with the
// message and the stack trace.
package preview

import (
	_ "embed"
	"fmt"
	"html"
	"runtime/debug"
	"time"

	gomdgo "github.com/polestorm/gomdgo/markdown"
	"github.com/polestorm/gomermaidgo/markdown"
	"github.com/polestorm/gomermaidgo/mermaid"
)

//go:embed style.css
var style string

// Page returns a complete HTML document that previews the Markdown source src, with title in <title>. The
// page has no scripts, its stylesheet and diagrams are inline, and raw HTML from the source is shown as text;
// only images the document itself links to by address are left for the viewer to load.
func Page(src []byte, title string) (out []byte) {
	defer func() {
		if r := recover(); r != nil {
			out = failure(title, fmt.Sprintf("panic: %v", r), string(debug.Stack()))
		}
	}()

	cfg := gomdgo.Config{
		Enable:     gomdgo.Alerts | gomdgo.Footnotes,
		HeadingIDs: true,
		CodeBlocks: map[string]gomdgo.CodeBlockFunc{"mermaid": diagram},
	}
	addHighlighting(cfg.CodeBlocks)
	start := time.Now()
	body, err := gomdgo.New(cfg).Render(src)
	elapsed := time.Since(start)
	if err != nil {
		return failure(title, err.Error(), "")
	}

	out = head(title)
	out = append(out, "<main>\n"...)
	out = append(out, body...)
	out = append(out, "</main>\n"...)
	out = fmt.Appendf(out, "<footer>gopeekgo %s · %s · %d B → %d B · %s</footer>\n",
		Version, libraries, len(src), len(body), elapsed.Round(10*time.Microsecond))
	return append(out, "</body>\n</html>\n"...)
}

// diagram draws one mermaid block. An error does not stop the preview: the message with the line number
// and the source of the block take the place of the diagram.
func diagram(b gomdgo.CodeBlock) (out []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			out = diagramError(b, fmt.Sprintf("panic: %v", r))
			err = nil
		}
	}()
	svg, renderErr := mermaid.Render(string(b.Code))
	if renderErr != nil {
		return diagramError(b, renderErr.Error()), nil
	}
	out = append([]byte("<figure class=\"diagram\">\n"), svg...)
	if n := len(out); n > 0 && out[n-1] != '\n' {
		out = append(out, '\n')
	}
	return append(out, "</figure>\n"...), nil
}

func diagramError(b gomdgo.CodeBlock, message string) []byte {
	return []byte(fmt.Sprintf("<pre class=\"diagram-error\"><strong>mermaid, line %d: %s</strong>%s</pre>\n",
		b.Line, html.EscapeString(message), html.EscapeString(string(b.Code))))
}

func head(title string) []byte {
	return []byte("<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n<title>" +
		html.EscapeString(title) + "</title>\n<style>" + markdown.CSS + style + highlightCSS + "</style>\n</head>\n<body>\n")
}

func failure(title, message, detail string) []byte {
	out := head(title)
	out = append(out, "<main>\n<pre class=\"diagram-error\"><strong>gomdgo: "...)
	out = append(out, html.EscapeString(message)...)
	out = append(out, "</strong>"...)
	out = append(out, html.EscapeString(detail)...)
	return append(out, "</pre>\n</main>\n</body>\n</html>\n"...)
}
