//go:build nohighlight

package preview

import gomdgo "github.com/polestorm/gomdgo/markdown"

const (
	libraries    = "gomdgo + gomermaidgo"
	highlightCSS = ""
)

// addHighlighting does nothing with the nohighlight tag: gomdgo writes code blocks as plain text.
func addHighlighting(map[string]gomdgo.CodeBlockFunc) {}
