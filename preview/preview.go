// Package preview składa stronę HTML podglądu pliku Markdown: gomdgo zamienia dokument na HTML, a bloki
// kodu o języku "mermaid" rysuje gomermaidgo. Pakiet nie zależy od cgo ani od macOS, więc da się go
// testować wszędzie; rozszerzenie Quick Look woła go przez ../bridge.
//
// Podgląd służy do oglądania dowolnych plików, także niepoprawnych, dlatego Page nie zwraca błędu:
// błąd diagramu trafia na stronę w miejscu diagramu, a panika biblioteki zamienia się w stronę
// z komunikatem (i jest sygnałem, że znaleźliśmy błąd do zgłoszenia).
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

// Page zwraca kompletny dokument HTML z podglądem źródła src. Tytuł trafia do <title>.
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
	out = fmt.Appendf(out, "<footer>gopeekgo · gomdgo + gomermaidgo · %d B → %d B · %s</footer>\n",
		len(src), len(body), elapsed.Round(10*time.Microsecond))
	return append(out, "</body>\n</html>\n"...)
}

// diagram rysuje jeden blok mermaid. Błąd nie przerywa podglądu: w miejscu diagramu zostaje
// komunikat z numerem wiersza i źródło bloku.
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
		html.EscapeString(title) + "</title>\n<style>" + markdown.CSS + style + "</style>\n</head>\n<body>\n")
}

func failure(title, message, detail string) []byte {
	out := head(title)
	out = append(out, "<main>\n<pre class=\"diagram-error\"><strong>gomdgo: "...)
	out = append(out, html.EscapeString(message)...)
	out = append(out, "</strong>"...)
	out = append(out, html.EscapeString(detail)...)
	return append(out, "</pre>\n</main>\n</body>\n</html>\n"...)
}
