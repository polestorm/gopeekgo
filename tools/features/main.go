// Command features repeats the feature rows of the comparison in the README. It converts each frozen probe
// file in testdata/features with gopeekgo (cmd/gopeekgo-render) and with another Markdown preview run as a
// black box, and reports what appears in each page by fixed rules written below. It prints the environment
// and a Markdown table, ready for docs/testing.md. The method is described in docs/benchmark.md.
//
// Run it from the root of the repository:
//
//	go run ./tools/features
//	go run ./tools/features -rival /path/to/qlmarkdown_cli
//	go run ./tools/features -tags nohighlight
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// A probe is one feature: a file in testdata/features and the rule that reads the result from the page.
// The rules read the content of the page without its <style> and <script> elements, so that class names in
// a stylesheet or a script are not taken for content; the Math rule also looks at the scripts.
type probe struct {
	feature string
	file    string
	read    func(p page) string
}

type page struct {
	full    string // the whole page
	content string // the page without <style> and <script> elements
}

var styleOrScript = regexp.MustCompile(`(?is)<style[^>]*>.*?</style>|<script[^>]*>.*?</script>`)

func newPage(full string) page { return page{full, styleOrScript.ReplaceAllString(full, "")} }

// has returns a rule that reports yes when the page contains any of the given strings, otherwise no.
func has(subs ...string) func(page) string {
	return func(p page) string {
		for _, s := range subs {
			if strings.Contains(p.content, s) {
				return "yes"
			}
		}
		return "no"
	}
}

var scriptSrc = regexp.MustCompile(`<script[^>]*\ssrc=["']([^"']+)["']`)

// hosts lists the hosts of the scripts the page loads by address.
func hosts(page string) []string {
	seen := map[string]bool{}
	for _, m := range scriptSrc.FindAllStringSubmatch(page, -1) {
		if host, ok := strings.CutPrefix(m[1], "https://"); ok {
			host, _, _ = strings.Cut(host, "/")
			seen[host] = true
		}
	}
	var out []string
	for h := range seen {
		out = append(out, h)
	}
	sort.Strings(out)
	return out
}

var probes = []probe{
	{"Tables", "table.md", has("<table>")},
	{"Task lists", "tasklist.md", has(`type="checkbox"`)},
	{"Strikethrough", "strikethrough.md", has("<del>struck</del>")},
	{"Bare `www.` links", "autolink.md", has(`href="http://www.example.com"`)},
	{"Footnotes", "footnote.md", has("footnote")},
	{"Heading identifiers", "heading-id.md", has(`<h1 id="heading-one"`)},
	{"GitHub alerts (`> [!NOTE]`)", "alert.md", has("alert-note")},
	{"Mermaid diagrams", "mermaid.md", func(p page) string {
		page := p.content
		switch {
		case strings.Contains(page, "<svg"):
			return "SVG in the page"
		case strings.Contains(page, `class="mermaid"`):
			return "Mermaid.js in the preview"
		}
		return "no"
	}},
	{"Syntax highlighting", "code.md", has(`<span class="hl`)},
	{"Math", "math.md", func(p page) string {
		if strings.Contains(p.full, "MathJax") || strings.Contains(p.full, "katex") {
			return "yes, a script in the preview"
		}
		return "no"
	}},
	{"YAML front matter", "front-matter.md", func(p page) string {
		page := p.content
		switch {
		case strings.Contains(page, "yaml"):
			return "shown as metadata"
		case strings.Contains(page, "<h2") && strings.Contains(page, "title: front"):
			return "read as Markdown"
		}
		return "no"
	}},
	{"Emoji shortcodes (`:smile:`)", "emoji.md", has("😄")},
	{"Image with a relative path", "image.md", func(p page) string {
		page := p.content
		switch {
		case strings.Contains(page, `src="data:image/png`):
			return "embedded in the page"
		case strings.Contains(page, `src="image.png"`):
			return "left as an address"
		}
		return "no"
	}},
	{"`<script>` in the source", "raw-script.md", func(p page) string {
		page := p.content
		switch {
		case strings.Contains(page, "<script>alert(1)"):
			return "kept as a script"
		case strings.Contains(page, "alert(1)"):
			return "shown as text"
		}
		return "removed"
	}},
	{"Other raw HTML (`<b>`)", "raw-tag.md", func(p page) string {
		page := p.content
		if strings.Contains(page, "<b>bold</b>") {
			return "rendered"
		}
		return "shown as text"
	}},
	{"`javascript:` link", "javascript-link.md", func(p page) string {
		page := p.content
		if strings.Contains(page, `href="javascript:`) {
			return "kept"
		}
		return "removed"
	}},
}

func main() {
	var (
		rival    = flag.String("rival", "/Applications/QLMarkdown.app/Contents/Resources/qlmarkdown_cli", "command-line converter of the other project")
		rivalApp = flag.String("rival-app", "/Applications/QLMarkdown.app", "installed application of the other project")
		tags     = flag.String("tags", "", "build tags for gopeekgo-render, for example nohighlight")
	)
	flag.Parse()

	work, err := os.MkdirTemp("", "gopeekgo-features-")
	check(err)
	defer os.RemoveAll(work)
	own := filepath.Join(work, "gopeekgo-render")
	build := exec.Command("go", "build", "-trimpath", "-tags", *tags, "-o", own, "./cmd/gopeekgo-render")
	build.Stderr = os.Stderr
	check(build.Run())

	_, err = os.Stat(*rival)
	rivalFound := err == nil

	fmt.Println("## Environment")
	fmt.Println()
	fmt.Printf("- date: %s\n", time.Now().Format("2006-01-02"))
	fmt.Printf("- gopeekgo: %s, commit %s, build tags %q\n", output(own, "-version"),
		output("git", "rev-parse", "--short", "HEAD"), *tags)
	if rivalFound {
		fmt.Printf("- other project: %s, version %s\n", *rival,
			output("plutil", "-extract", "CFBundleShortVersionString", "raw", filepath.Join(*rivalApp, "Contents", "Info.plist")))
		if settings := output(*rival, "--show-settings"); settings != "unknown" {
			fmt.Println("- settings of the other project (`--show-settings`):")
			fmt.Println()
			fmt.Println("```")
			fmt.Println(settings)
			fmt.Println("```")
		}
	} else {
		fmt.Printf("- other project: not found at %s (use -rival); only gopeekgo is probed\n", *rival)
	}
	fmt.Println()

	fmt.Println("## Features")
	fmt.Println()
	fmt.Println("| feature | probe file | gopeekgo | other |")
	fmt.Println("|---|---|---|---|")
	var ownHosts, otherHosts []string
	ownScripts, otherScripts := 0, 0
	for _, p := range probes {
		path := filepath.Join("testdata", "features", p.file)
		a := convert(own, path)
		ownScripts += strings.Count(a, "<script")
		ownHosts = append(ownHosts, hosts(a)...)
		b := "—"
		if rivalFound {
			other := convert(*rival, path)
			otherScripts += strings.Count(other, "<script")
			otherHosts = append(otherHosts, hosts(other)...)
			b = p.read(newPage(other))
		}
		fmt.Printf("| %s | `%s` | %s | %s |\n", p.feature, p.file, p.read(newPage(a)), b)
	}
	other := "—"
	if rivalFound {
		other = summary(otherScripts, otherHosts)
	}
	fmt.Printf("| Scripts in the pages, hosts they load from | all files | %s | %s |\n", summary(ownScripts, ownHosts), other)
}

// convert runs the program on the file, from the file's directory, so a relative image can be found.
func convert(program, path string) string {
	abs, err := filepath.Abs(program)
	check(err)
	cmd := exec.Command(abs, filepath.Base(path))
	cmd.Dir = filepath.Dir(path)
	out, err := cmd.Output()
	if err != nil {
		check(fmt.Errorf("%s %s: %w", program, path, err))
	}
	return string(out)
}

func summary(scripts int, hostList []string) string {
	if scripts == 0 {
		return "none"
	}
	seen := map[string]bool{}
	var unique []string
	for _, h := range hostList {
		if !seen[h] {
			seen[h] = true
			unique = append(unique, h)
		}
	}
	sort.Strings(unique)
	if len(unique) == 0 {
		return fmt.Sprintf("%d `<script>` elements, no hosts", scripts)
	}
	return fmt.Sprintf("%d `<script>` elements; loaded from %s", scripts, strings.Join(unique, ", "))
}

func output(name string, args ...string) string {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "features:", err)
		os.Exit(1)
	}
}
