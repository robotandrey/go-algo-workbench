package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/robotandrey/go-algo-workbench/internal/leetcode"
)

const usage = `fetch — scaffold a LeetCode problem into go-algo-workbench

Usage:
  fetch -url <leetcode-url>
  fetch -slug <problem-slug>

Examples:
  fetch -url https://leetcode.com/problems/two-sum/
  fetch -slug valid-parentheses

Flags:
`

func main() {
	var (
		url  = flag.String("url", "", "Full LeetCode problem URL")
		slug = flag.String("slug", "", "Problem slug (e.g. two-sum)")
	)
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
	}
	flag.Parse()

	// Also accept positional arg: fetch two-sum
	if *url == "" && *slug == "" && flag.NArg() > 0 {
		arg := flag.Arg(0)
		if len(arg) > 20 && arg[:8] == "https://" {
			*url = arg
		} else {
			*slug = arg
		}
	}

	if *url == "" && *slug == "" {
		flag.Usage()
		os.Exit(1)
	}

	input := *url
	if input == "" {
		input = *slug
	}

	fmt.Printf("→ Fetching problem: %s\n", input)

	problem, err := leetcode.FetchProblem(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Got: [%s] %s (%s)\n", problem.FrontendID, problem.Title, problem.Difficulty)

	problemsDir := findProblemsDir()
	dir, err := leetcode.Scaffold(problem, problemsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Scaffolded: %s\n", dir)
	fmt.Println()
	fmt.Printf("  Files created:\n")
	fmt.Printf("    %s/problem.md       ← problem statement\n", filepath.Base(dir))
	fmt.Printf("    %s/solution.go      ← your solution\n", filepath.Base(dir))
	fmt.Printf("    %s/solution_test.go ← table-driven tests\n", filepath.Base(dir))
	fmt.Println()

	// Print run hint
	dirName := filepath.Base(dir)
	if runtime.GOOS == "windows" {
		fmt.Printf("  Run tests:  make one name=%s\n", dirName)
	} else {
		fmt.Printf("  Run tests:  make one name=%s\n", dirName)
	}
}

// findProblemsDir walks up from the binary location to find the problems/ directory.
// Falls back to ./problems relative to cwd.
func findProblemsDir() string {
	// Try cwd first (when running via `go run` or `make`)
	cwd, err := os.Getwd()
	if err == nil {
		candidate := filepath.Join(cwd, "problems")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}

	// Try executable location
	exe, err := os.Executable()
	if err == nil {
		for dir := filepath.Dir(exe); dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
			candidate := filepath.Join(dir, "problems")
			if info, err := os.Stat(candidate); err == nil && info.IsDir() {
				return candidate
			}
		}
	}

	// Fallback: create problems/ in cwd
	return filepath.Join(cwd, "problems")
}
