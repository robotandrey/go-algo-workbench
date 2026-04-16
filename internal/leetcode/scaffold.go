package leetcode

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Scaffold creates the problem directory with solution.go, solution_test.go, and problem.md.
func Scaffold(p *Problem, problemsDir string) (string, error) {
	dirName := strings.ReplaceAll(p.Slug, "-", "_")
	dir := filepath.Join(problemsDir, dirName)

	if _, err := os.Stat(dir); err == nil {
		return "", fmt.Errorf("problem directory already exists: %s", dir)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	gen := &GoCodegen{Problem: p}

	files := map[string]string{
		"solution.go":      gen.GenerateSolution(),
		"solution_test.go": gen.GenerateTest(),
		"problem.md":       GenerateMarkdown(p),
	}

	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return "", fmt.Errorf("failed to write %s: %w", name, err)
		}
	}

	return dir, nil
}
