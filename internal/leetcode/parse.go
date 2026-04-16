package leetcode

import (
	"encoding/json"
	"html"
	"regexp"
	"strings"
)

// metaData is the structure of the metaData JSON field from LeetCode.
type metaData struct {
	Params []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"params"`
	Return struct {
		Type string `json:"type"`
	} `json:"return"`
}

var labeledLineRe = regexp.MustCompile(`^[A-Za-z][A-Za-z ]*:`)

// parseExampleTests pairs up raw test case lines with expected outputs.
// Inputs are read from exampleTestcaseList and outputs are extracted
// from HTML example blocks in content.
func parseExampleTests(inputs []string, rawMeta, content string) []SampleTest {
	var meta metaData
	_ = json.Unmarshal([]byte(rawMeta), &meta)

	outputs := extractExampleOutputs(content)

	var tests []SampleTest
	for i, raw := range inputs {
		lines := splitInputLines(raw)
		var parts []string
		var args []string
		for idx, line := range lines {
			value := line
			args = append(args, line)
			if idx < len(meta.Params) {
				name := strings.TrimSpace(meta.Params[idx].Name)
				if name != "" {
					value = stripNamedValuePrefix(line, name)
				}
				args[len(args)-1] = value
				if name != "" {
					parts = append(parts, name+" = "+value)
				} else {
					parts = append(parts, value)
				}
			} else {
				parts = append(parts, value)
			}
		}

		output := "TODO"
		if i < len(outputs) && outputs[i] != "" {
			output = outputs[i]
		}

		tests = append(tests, SampleTest{
			Input:  strings.Join(parts, "\n"),
			Output: output,
			Args:   args,
		})
	}
	return tests
}

func splitInputLines(raw string) []string {
	var out []string
	for _, line := range strings.Split(strings.TrimSpace(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}

func extractExampleOutputs(content string) []string {
	outputs := extractOutputsFromPreBlocks(content)
	if len(outputs) > 0 {
		return outputs
	}
	return extractOutputsFromParagraphBlocks(content)
}

func extractOutputsFromPreBlocks(content string) []string {
	rePre := regexp.MustCompile(`(?is)<pre[^>]*>(.*?)</pre>`)
	matches := rePre.FindAllStringSubmatch(content, -1)

	var out []string
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		block := html.UnescapeString(stripTags(m[1]))
		output := extractOutputFromExampleBlock(block)
		if output != "" {
			out = append(out, output)
		}
	}
	return out
}

func extractOutputsFromParagraphBlocks(content string) []string {
	reOutput := regexp.MustCompile(`(?is)<strong[^>]*>\s*Output:\s*</strong>\s*(?:<span[^>]*>)?\s*(.*?)\s*(?:</span>)?\s*</p>`)
	matches := reOutput.FindAllStringSubmatch(content, -1)

	var out []string
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		value := strings.TrimSpace(stripTags(html.UnescapeString(m[1])))
		if value != "" {
			out = append(out, value)
		}
	}
	return out
}

func extractOutputFromExampleBlock(block string) string {
	lines := strings.Split(strings.ReplaceAll(block, "\r\n", "\n"), "\n")
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(line, "Output:") {
			continue
		}

		value := strings.TrimSpace(strings.TrimPrefix(line, "Output:"))
		if value != "" {
			return value
		}

		var multi []string
		for j := i + 1; j < len(lines); j++ {
			next := strings.TrimSpace(lines[j])
			if next == "" {
				if len(multi) > 0 {
					break
				}
				continue
			}
			if labeledLineRe.MatchString(next) {
				break
			}
			multi = append(multi, next)
		}
		return strings.Join(multi, " ")
	}
	return ""
}

func stripNamedValuePrefix(line, paramName string) string {
	left, right, ok := strings.Cut(line, "=")
	if !ok {
		return strings.TrimSpace(line)
	}
	if strings.TrimSpace(left) == paramName {
		return strings.TrimSpace(right)
	}
	return strings.TrimSpace(line)
}
