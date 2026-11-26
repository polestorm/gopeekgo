// Package rival holds what tools/bench and tools/features need to run the compared program the same way on
// every machine.
package rival

import "path/filepath"

// QLMarkdownFlags fixes every setting of qlmarkdown_cli that has a command-line flag, so the result does not
// depend on the settings saved in the user's QLMarkdown application. The values are the settings of a
// QLMarkdown 1.5.7 installation on 2026-11-26 with one change: GitHub alerts are on, because gopeekgo
// renders them too. The appearance is fixed to light. The flags not listed here (the base font size) have no
// effect on the features compared.
var QLMarkdownFlags = []string{
	"--appearance", "light",
	"--footnotes", "on",
	"--hard-break", "off",
	"--no-soft-break", "off",
	"--raw-html", "on",
	"--render-as-code", "off",
	"--smart-quotes", "on",
	"--table-of-contents", "off",
	"--validate-utf8", "off",
	"--debug", "off",
	"--admonition", "off",
	"--alert", "on",
	"--autolink", "on",
	"--definition-list", "off",
	"--emoji", "font",
	"--github-mentions", "off",
	"--hashtag", "off",
	"--heads-anchor", "on",
	"--highlight", "off",
	"--inline-images", "on",
	"--math", "link",
	"--mermaid", "link",
	"--strikethrough", "single",
	"--sub", "off",
	"--sup", "off",
	"--syntax-highlight", "on",
	"--table", "on",
	"--tag-filter", "on",
	"--tasklist", "on",
	"--wikilink", "off",
	"--yaml", "all",
}

// Args returns the arguments for converting file with program: QLMarkdownFlags first when program is
// QLMarkdown's converter, then the file name. Another converter gets the file name only.
func Args(program, file string) []string {
	if filepath.Base(program) == "qlmarkdown_cli" {
		return append(append([]string{}, QLMarkdownFlags...), file)
	}
	return []string{file}
}
