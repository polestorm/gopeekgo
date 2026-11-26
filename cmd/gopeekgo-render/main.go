// Command gopeekgo-render writes to standard output exactly the HTML page that the gopeekgo Quick Look
// extension returns for a Markdown file (preview.Page): CommonMark and GFM, Mermaid diagrams as SVG and
// highlighted code blocks in one self-contained page. It is used to look at the result outside Finder and
// to measure it (docs/benchmark.md).
//
// Usage:
//
//	gopeekgo-render file.md > page.html
//	gopeekgo-render < file.md > page.html
//	gopeekgo-render -version
//
// Building with -tags nohighlight leaves syntax highlighting out, as for the extension.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/polestorm/gopeekgo/preview"
)

func main() {
	version := flag.Bool("version", false, "print the version and exit")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: gopeekgo-render [-version] [file.md]")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *version {
		fmt.Println("gopeekgo-render", preview.Version)
		return
	}
	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(2)
	}
	var (
		src   []byte
		err   error
		title = "stdin"
	)
	if flag.NArg() == 1 {
		title = filepath.Base(flag.Arg(0))
		src, err = os.ReadFile(flag.Arg(0))
	} else {
		src, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "gopeekgo-render:", err)
		os.Exit(1)
	}
	if _, err := os.Stdout.Write(preview.Page(src, title)); err != nil {
		fmt.Fprintln(os.Stderr, "gopeekgo-render:", err)
		os.Exit(1)
	}
}
