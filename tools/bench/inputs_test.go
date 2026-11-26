package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInputsAreFrozen checks that the input of the comparison is byte for byte the input listed in
// testdata/SHA256SUMS, and that no file was added to testdata/bench or testdata/features without a sum.
// Changing the input is changing the benchmark: update SHA256SUMS and docs/benchmark.md and measure again.
func TestInputsAreFrozen(t *testing.T) {
	root := filepath.Join("..", "..", "testdata")
	list, err := os.Open(filepath.Join(root, "SHA256SUMS"))
	if err != nil {
		t.Fatal(err)
	}
	defer list.Close()

	want := map[string]string{}
	lines := bufio.NewScanner(list)
	for lines.Scan() {
		sum, name, ok := strings.Cut(lines.Text(), "  ")
		if !ok {
			t.Fatalf("malformed line in SHA256SUMS: %q", lines.Text())
		}
		want[name] = sum
	}
	if err := lines.Err(); err != nil {
		t.Fatal(err)
	}

	for _, dir := range []string{"bench", "features"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			name, _ := filepath.Rel(root, path)
			name = filepath.ToSlash(name)
			if _, listed := want[name]; !listed {
				t.Errorf("%s is not listed in testdata/SHA256SUMS", name)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	for name, sum := range want {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		got := sha256.Sum256(data)
		if hex.EncodeToString(got[:]) != sum {
			t.Errorf("%s changed: SHA-256 %x, listed %s", name, got, sum)
		}
	}
}
