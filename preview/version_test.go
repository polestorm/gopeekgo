package preview

import (
	"os"
	"regexp"
	"testing"
)

func TestVersionMatchesBundles(t *testing.T) {
	re := regexp.MustCompile(`<key>CFBundleShortVersionString</key>\s*<string>([^<]*)</string>`)
	for _, name := range []string{"../macos/App-Info.plist", "../macos/Extension-Info.plist"} {
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		m := re.FindSubmatch(data)
		if m == nil || string(m[1]) != Version {
			t.Errorf("%s: CFBundleShortVersionString is %q, preview.Version is %q", name, m, Version)
		}
	}
}
