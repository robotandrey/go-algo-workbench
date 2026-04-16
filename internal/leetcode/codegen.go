package leetcode

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// GoTypeMap maps LeetCode type strings to Go types.
var GoTypeMap = map[string]string{
	"integer":             "int",
	"long":                "int64",
	"float":               "float64",
	"double":              "float64",
	"string":              "string",
	"boolean":             "bool",
	"character":           "byte",
	"integer[]":           "[]int",
	"long[]":              "[]int64",
	"double[]":            "[]float64",
	"string[]":            "[]string",
	"boolean[]":           "[]bool",
	"character[]":         "[]byte",
	"integer[][]":         "[][]int",
	"string[][]":          "[][]string",
	"list<integer>":       "[]int",
	"list<string>":        "[]string",
	"list<double>":        "[]float64",
	"list<list<integer>>": "[][]int",
	"void":                "",
}

// mapGoType converts a LeetCode type string to a Go type.
func mapGoType(lcType string) string {
	lower := strings.ToLower(lcType)
	if goType, ok := GoTypeMap[lower]; ok {
		return goType
	}
	// Fallback: return as-is (e.g. TreeNode, ListNode)
	return lcType
}

// GoCodegen generates solution.go and solution_test.go content
// from the fetched problem and its Go snippet.
type GoCodegen struct {
	Problem *Problem
}

// GenerateSolution returns the content of solution.go.
// If the problem has a Go snippet from LeetCode, it uses that as a base.
func (g *GoCodegen) GenerateSolution() string {
	pkg := PackageName(g.Problem.Slug)
	if g.Problem.CodeSnippet != "" {
		snippet := cleanSnippet(g.Problem.CodeSnippet)
		return fmt.Sprintf("package %s\n\n%s\n", pkg, snippet)
	}
	// Fallback to generic template
	return fmt.Sprintf("package %s\n\n// Solve is the entry point for the solution.\nfunc Solve() {\n}\n", pkg)
}

// GenerateTest returns the content of solution_test.go.
func (g *GoCodegen) GenerateTest() string {
	pkg := PackageName(g.Problem.Slug)
	funcName, params, retType := parseFuncFromSnippet(g.Problem.CodeSnippet)

	if funcName == "" {
		return g.genericTest(pkg)
	}

	return g.typedTest(pkg, funcName, params, retType)
}

type param struct {
	name   string
	goType string
}

// parseFuncFromSnippet extracts function name, params and return type from a Go snippet.
func parseFuncFromSnippet(snippet string) (name string, params []param, retType string) {
	// Match: func FuncName(arg1 type1, arg2 type2) retType {
	re := regexp.MustCompile(`func\s+([A-Za-z][A-Za-z0-9]*)\s*\(([^)]*)\)\s*([^{]*)`)
	m := re.FindStringSubmatch(snippet)
	if m == nil {
		return "", nil, ""
	}
	name = m[1]
	retType = strings.TrimSpace(m[3])

	// Parse params
	rawParams := strings.TrimSpace(m[2])
	if rawParams != "" {
		for _, part := range strings.Split(rawParams, ",") {
			part = strings.TrimSpace(part)
			fields := strings.Fields(part)
			if len(fields) >= 2 {
				params = append(params, param{
					name:   fields[0],
					goType: strings.Join(fields[1:], " "),
				})
			}
		}
	}
	return name, params, retType
}

// typedTest generates a properly typed table-driven test.
func (g *GoCodegen) typedTest(pkg, funcName string, params []param, retType string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("package %s\n\n", pkg))
	sb.WriteString("import \"testing\"\n\n")

	// Build struct fields
	sb.WriteString(fmt.Sprintf("func Test%s(t *testing.T) {\n", funcName))
	sb.WriteString("\ttests := []struct {\n")
	sb.WriteString("\t\tname string\n")
	for _, p := range params {
		sb.WriteString(fmt.Sprintf("\t\t%s %s\n", p.name, p.goType))
	}
	if retType != "" {
		sb.WriteString(fmt.Sprintf("\t\twant %s\n", retType))
	}
	sb.WriteString("\t}{\n")

	// Sample test cases from examples
	for i, tc := range g.Problem.SampleTests {
		caseName := fmt.Sprintf("example_%d", i+1)
		sb.WriteString(fmt.Sprintf("\t\t{\n\t\t\tname: %q,\n", caseName))
		sb.WriteString(fmt.Sprintf("\t\t\t// Input:  %s\n", strings.ReplaceAll(tc.Input, "\n", " | ")))
		if retType != "" {
			sb.WriteString(fmt.Sprintf("\t\t\t// Output: %s\n", tc.Output))
		}
		for idx, p := range params {
			value := zeroValue(p.goType)
			if idx < len(tc.Args) {
				if literal, ok := toGoLiteral(tc.Args[idx], p.goType); ok {
					value = literal
				}
			}
			sb.WriteString(fmt.Sprintf("\t\t\t%s: %s,\n", p.name, value))
		}
		if retType != "" {
			want := zeroValue(retType)
			if literal, ok := toGoLiteral(tc.Output, retType); ok {
				want = literal
			}
			sb.WriteString(fmt.Sprintf("\t\t\twant: %s,\n", want))
		}
		sb.WriteString("\t\t},\n")
	}
	if len(g.Problem.SampleTests) == 0 {
		sb.WriteString("\t\t{name: \"basic\"},\n")
	}

	sb.WriteString("\t}\n\n")
	sb.WriteString("\tfor _, tt := range tests {\n")
	sb.WriteString("\t\tt.Run(tt.name, func(t *testing.T) {\n")

	// Function call
	var argNames []string
	for _, p := range params {
		argNames = append(argNames, "tt."+p.name)
	}
	call := fmt.Sprintf("%s(%s)", funcName, strings.Join(argNames, ", "))
	if retType != "" {
		sb.WriteString(fmt.Sprintf("\t\t\tgot := %s\n", call))
		sb.WriteString("\t\t\t_ = got // TODO: compare got with tt.want\n")
		sb.WriteString("\t\t\t// if got != tt.want {\n")
		sb.WriteString("\t\t\t// \tt.Fatalf(\"got %%v, want %%v\", got, tt.want)\n")
		sb.WriteString("\t\t\t// }\n")
	} else {
		sb.WriteString(fmt.Sprintf("\t\t\t%s\n", call))
	}

	sb.WriteString("\t\t})\n")
	sb.WriteString("\t}\n")
	sb.WriteString("}\n")

	return sb.String()
}

func (g *GoCodegen) genericTest(pkg string) string {
	return fmt.Sprintf(`package %s

import "testing"

func TestSolve(t *testing.T) {
	tests := []struct {
		name string
		// TODO: add input/output fields
	}{
		{name: "basic"},
		{name: "edge_case"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// got := Solve(...)
			// if got != tt.want {
			// 	t.Fatalf("got %%v, want %%v", got, tt.want)
			// }
		})
	}
}
`, pkg)
}

// PackageName converts a slug like "two-sum" to a valid Go package name "twosum".
func PackageName(slug string) string {
	// Replace hyphens and underscores, keep only valid chars
	var sb strings.Builder
	for _, r := range strings.ToLower(slug) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			sb.WriteRune(r)
		}
	}
	name := sb.String()
	// Go package names cannot start with a digit
	if len(name) > 0 && unicode.IsDigit(rune(name[0])) {
		name = "p" + name
	}
	return name
}

// cleanSnippet removes comment markers LeetCode sometimes adds.
func cleanSnippet(snippet string) string {
	lines := strings.Split(snippet, "\n")
	var out []string
	for _, l := range lines {
		// Remove LeetCode's class Solution wrapper hints if present (Python-style, shouldn't be in Go)
		out = append(out, l)
	}
	return strings.Join(out, "\n")
}

// zeroValue returns a Go zero-value literal for a given type.
func zeroValue(t string) string {
	t = strings.TrimSpace(t)
	switch t {
	case "int", "int64", "float64":
		return "0"
	case "bool":
		return "false"
	case "string":
		return `""`
	case "byte":
		return "0"
	default:
		if strings.HasPrefix(t, "[]") || strings.HasPrefix(t, "map") {
			return "nil"
		}
		return fmt.Sprintf("%s{}", t)
	}
}

func toGoLiteral(raw, goType string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "TODO" {
		return "", false
	}

	switch raw {
	case "null":
		if strings.HasPrefix(goType, "[]") || strings.HasPrefix(goType, "map[") || strings.HasPrefix(goType, "*") {
			return "nil", true
		}
		return "", false
	}

	switch strings.TrimSpace(goType) {
	case "int", "int64", "float64", "bool":
		return raw, true
	case "string":
		if strings.HasPrefix(raw, `"`) && strings.HasSuffix(raw, `"`) {
			return raw, true
		}
		return fmt.Sprintf("%q", raw), true
	case "byte":
		if strings.HasPrefix(raw, `"`) && strings.HasSuffix(raw, `"`) && len(raw) == 3 {
			return "'" + raw[1:2] + "'", true
		}
		if strings.HasPrefix(raw, `'`) && strings.HasSuffix(raw, `'`) {
			return raw, true
		}
		return "", false
	}

	if strings.HasPrefix(goType, "[]") || strings.HasPrefix(goType, "[") {
		if strings.HasPrefix(raw, "[") && strings.HasSuffix(raw, "]") {
			body, ok := jsonArrayToGoBraces(raw)
			if !ok {
				return "", false
			}
			return goType + body, true
		}
	}
	if strings.HasPrefix(goType, "map[") {
		if strings.HasPrefix(raw, "{") && strings.HasSuffix(raw, "}") {
			return raw, true
		}
	}

	return "", false
}

func jsonArrayToGoBraces(raw string) (string, bool) {
	var b strings.Builder
	inString := false
	escaped := false

	for _, r := range raw {
		switch {
		case escaped:
			b.WriteRune(r)
			escaped = false
		case r == '\\':
			b.WriteRune(r)
			escaped = true
		case r == '"':
			b.WriteRune(r)
			inString = !inString
		case !inString && r == '[':
			b.WriteRune('{')
		case !inString && r == ']':
			b.WriteRune('}')
		default:
			b.WriteRune(r)
		}
	}

	if inString || escaped {
		return "", false
	}
	return b.String(), true
}
