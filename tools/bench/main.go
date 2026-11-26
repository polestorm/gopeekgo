// Command bench repeats the comparison of the README: conversion time and peak memory, and size on disk,
// for gopeekgo and for another Markdown preview run as a black box. It prints the environment and Markdown
// tables, ready for docs/testing.md. The input and the method are described in docs/benchmark.md.
//
// Run it from the root of the repository:
//
//	go run ./tools/bench
//	go run ./tools/bench -runs 30 -rival /path/to/qlmarkdown_cli
//	go run ./tools/bench -tags nohighlight      (gopeekgo without syntax highlighting)
package main

import (
	"bytes"
	"debug/macho"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/polestorm/gopeekgo/tools/internal/rival"
)

// The large inputs are copies of the frozen files joined together, about 1.3 MB each.
const (
	largeCopies     = 55 // copies of testdata/bench/document.md
	largeCodeCopies = 87 // copies of testdata/bench/code.md
)

type input struct {
	label string
	path  string
}

type result struct {
	median time.Duration
	min    time.Duration
	memory int64 // median peak memory, in bytes
	output int   // size of the output of the last run
}

func main() {
	var (
		runs      = flag.Int("runs", 15, "number of runs per program and input")
		rivalPath = flag.String("rival", "/Applications/QLMarkdown.app/Contents/Resources/qlmarkdown_cli", "command-line converter of the other project")
		rivalApp  = flag.String("rival-app", "/Applications/QLMarkdown.app", "installed application of the other project")
		ownApp    = flag.String("app", "", "gopeekgo application to measure; empty: build it with macos/build.sh and the same -tags")
		tags      = flag.String("tags", "", "build tags for gopeekgo-render, for example nohighlight")
	)
	flag.Parse()

	work, err := os.MkdirTemp("", "gopeekgo-bench-")
	check(err)
	defer os.RemoveAll(work)

	own := filepath.Join(work, "gopeekgo-render")
	build := exec.Command("go", "build", "-trimpath", "-tags", *tags, "-o", own, "./cmd/gopeekgo-render")
	build.Stderr = os.Stderr
	check(build.Run())

	// The size rows measure an application built from this tree with the same tags, not whatever is installed.
	// It is removed afterwards: a second bundle with the extension's identifier confuses the system registry.
	if *ownApp == "" && runtime.GOOS == "darwin" {
		script := exec.Command("macos/build.sh")
		script.Env = append(os.Environ(), "TAGS="+*tags)
		script.Stdout, script.Stderr = os.Stderr, os.Stderr
		check(script.Run())
		*ownApp = filepath.Join("build", "GoPeekGo.app")
		defer os.RemoveAll(*ownApp)
	}

	document, err := os.ReadFile("testdata/bench/document.md")
	check(err)
	large := filepath.Join(work, "large.md")
	check(os.WriteFile(large, bytes.Repeat(document, largeCopies), 0o644))
	code, err := os.ReadFile("testdata/bench/code.md")
	check(err)
	largeCode := filepath.Join(work, "large-code.md")
	check(os.WriteFile(largeCode, bytes.Repeat(code, largeCodeCopies), 0o644))

	inputs := []input{
		{"document.md", "testdata/bench/document.md"},
		{"diagram.md", "testdata/bench/diagram.md"},
		{fmt.Sprintf("%d × document.md", largeCopies), large},
		{"code.md", "testdata/bench/code.md"},
		{fmt.Sprintf("%d × code.md", largeCodeCopies), largeCode},
	}

	rivalFound := true
	if _, err := os.Stat(*rivalPath); err != nil {
		rivalFound = false
	}

	fmt.Println("## Environment")
	fmt.Println()
	fmt.Printf("- date: %s\n", time.Now().Format("2006-01-02"))
	fmt.Printf("- machine: %s, %s/%s\n", output("sysctl", "-n", "machdep.cpu.brand_string"), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("- system: %s %s\n", output("sw_vers", "-productName"), output("sw_vers", "-productVersion"))
	fmt.Printf("- power: %s\n", firstLine(output("pmset", "-g", "batt")))
	fmt.Printf("- Go: %s\n", output("go", "env", "GOVERSION"))
	fmt.Printf("- gopeekgo: %s, commit %s%s, build tags %q\n", output(own, "-version"),
		output("git", "rev-parse", "--short", "HEAD"), dirty(), *tags)
	fmt.Printf("- libraries: %s\n", strings.ReplaceAll(output("go", "list", "-m", "github.com/polestorm/gomdgo",
		"github.com/polestorm/gomermaidgo", "github.com/polestorm/gohighlightgo"), "\n", ", "))
	if rivalFound {
		fmt.Printf("- other project: %s, version %s\n", *rivalPath,
			output("plutil", "-extract", "CFBundleShortVersionString", "raw", filepath.Join(*rivalApp, "Contents", "Info.plist")))
		fmt.Printf("- arguments of the other project: `%s`\n", strings.Join(rival.Args(*rivalPath, "FILE"), " "))
	} else {
		fmt.Printf("- other project: not found at %s (use -rival); only gopeekgo is measured\n", *rivalPath)
	}
	fmt.Printf("- runs per cell: %d\n", *runs)
	fmt.Println()

	fmt.Println("## Conversion: time and peak memory")
	fmt.Println()
	fmt.Println("| input | size | gopeekgo median (min) | other median (min) | advantage | gopeekgo memory | other memory | advantage |")
	fmt.Println("|---|--:|--:|--:|--:|--:|--:|--:|")
	for _, in := range inputs {
		info, err := os.Stat(in.path)
		check(err)
		a := measure(*runs, own, in.path)
		row := fmt.Sprintf("| %s | %d B | %s (%s) |", in.label, info.Size(), ms(a.median), ms(a.min))
		if rivalFound {
			b := measure(*runs, *rivalPath, in.path)
			row += fmt.Sprintf(" %s (%s) | %s | %s | %s | %s |", ms(b.median), ms(b.min),
				advantage(float64(a.median), float64(b.median)), mb(a.memory), mb(b.memory),
				advantage(float64(a.memory), float64(b.memory)))
		} else {
			row += fmt.Sprintf(" — | — | %s | — | — |", mb(a.memory))
		}
		fmt.Println(row)
	}
	fmt.Println()

	fmt.Println("## Size on disk")
	fmt.Println()
	fmt.Println("| | gopeekgo | other | advantage |")
	fmt.Println("|---|--:|--:|--:|")
	sizeRow("application bundle, sum of file sizes", treeSize(*ownApp), treeSize(*rivalApp))
	sizeRow("Quick Look extension bundle", treeSize(extension(*ownApp)), treeSize(extension(*rivalApp)))
	sizeRow("extension executable, "+runtime.GOARCH+" slice", sliceSize(extension(*ownApp)), sliceSize(extension(*rivalApp)))
}

// measure runs the program on the file the given number of times, after one warm-up run that is not
// counted, and returns the median wall-clock time of the whole process and the median peak memory.
func measure(runs int, program, path string) result {
	var times []time.Duration
	var memory []int64
	var size int
	for i := 0; i <= runs; i++ {
		var out bytes.Buffer
		cmd := exec.Command(program, rival.Args(program, path)...)
		cmd.Stdout = &out
		start := time.Now()
		err := cmd.Run()
		elapsed := time.Since(start)
		if err != nil {
			check(fmt.Errorf("%s %s: %w", program, path, err))
		}
		if i == 0 {
			continue
		}
		times = append(times, elapsed)
		memory = append(memory, peakMemory(cmd.ProcessState))
		size = out.Len()
	}
	sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
	sort.Slice(memory, func(i, j int) bool { return memory[i] < memory[j] })
	return result{median: times[len(times)/2], min: times[0], memory: memory[len(memory)/2], output: size}
}

func peakMemory(state *os.ProcessState) int64 {
	usage, ok := state.SysUsage().(*syscall.Rusage)
	if !ok {
		return 0
	}
	if runtime.GOOS == "darwin" {
		return int64(usage.Maxrss) // bytes
	}
	return int64(usage.Maxrss) * 1024 // kilobytes on Linux
}

// extension returns the path of the first extension in Contents/PlugIns of the application, or "".
func extension(app string) string {
	matches, _ := filepath.Glob(filepath.Join(app, "Contents", "PlugIns", "*.appex"))
	if len(matches) == 0 {
		return ""
	}
	return matches[0]
}

func treeSize(root string) int64 {
	if root == "" {
		return 0
	}
	var total int64
	filepath.WalkDir(root, func(_ string, d fs.DirEntry, err error) error {
		if err != nil || !d.Type().IsRegular() {
			return nil
		}
		if info, err := d.Info(); err == nil {
			total += info.Size()
		}
		return nil
	})
	return total
}

// sliceSize returns the size of the extension executable for the architecture of this machine: the whole
// file when it has one architecture, or the matching slice of a universal file.
func sliceSize(appex string) int64 {
	matches, _ := filepath.Glob(filepath.Join(appex, "Contents", "MacOS", "*"))
	if appex == "" || len(matches) == 0 {
		return 0
	}
	want := macho.CpuArm64
	if runtime.GOARCH == "amd64" {
		want = macho.CpuAmd64
	}
	if fat, err := macho.OpenFat(matches[0]); err == nil {
		defer fat.Close()
		for _, arch := range fat.Arches {
			if arch.Cpu == want {
				return int64(arch.Size)
			}
		}
		return 0
	}
	info, err := os.Stat(matches[0])
	if err != nil {
		return 0
	}
	return info.Size()
}

func sizeRow(label string, own, other int64) {
	cell := func(n int64) string {
		if n == 0 {
			return "not installed"
		}
		return fmt.Sprintf("%d B (%s)", n, mb(n))
	}
	adv := "—"
	if own > 0 && other > 0 {
		adv = advantage(float64(own), float64(other))
	}
	fmt.Printf("| %s | %s | %s | %s |\n", label, cell(own), cell(other), adv)
}

// advantage is gopeekgo's advantage: (other − gopeekgo) / other × 100%.
func advantage(own, other float64) string {
	if other == 0 {
		return "—"
	}
	return fmt.Sprintf("%.0f%%", (other-own)/other*100)
}

func ms(d time.Duration) string { return fmt.Sprintf("%.1f ms", float64(d.Microseconds())/1000) }

func mb(n int64) string { return fmt.Sprintf("%.1f MB", float64(n)/1e6) }

func output(name string, args ...string) string {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func firstLine(s string) string {
	line, _, _ := strings.Cut(s, "\n")
	return line
}

func dirty() string {
	if output("git", "status", "--porcelain") != "" {
		return " with uncommitted changes"
	}
	return ""
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "bench:", err)
		os.Exit(1)
	}
}
