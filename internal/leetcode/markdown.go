package leetcode

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"
)

// GenerateMarkdown creates the content for problem.md.
func GenerateMarkdown(p *Problem) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("# %s. %s\n\n", p.FrontendID, p.Title))
	sb.WriteString(fmt.Sprintf("**Difficulty:** %s  \n", p.Difficulty))
	sb.WriteString(fmt.Sprintf("**URL:** https://leetcode.com/problems/%s/  \n", p.Slug))
	sb.WriteString(fmt.Sprintf("**Fetched:** %s\n\n", time.Now().Format("2006-01-02")))

	if len(p.Topics) > 0 {
		var tagged []string
		for _, t := range p.Topics {
			tagged = append(tagged, "`"+t+"`")
		}
		sb.WriteString(fmt.Sprintf("**Topics:** %s\n\n", strings.Join(tagged, " ")))
	}

	sb.WriteString("---\n\n")

	// Problem statement (convert HTML → Markdown)
	if p.Content != "" {
		sb.WriteString("## Problem\n\n")
		sb.WriteString(htmlToMarkdown(p.Content))
		sb.WriteString("\n\n")
	}

	// Examples
	if len(p.SampleTests) > 0 {
		sb.WriteString("## Examples\n\n")
		for i, tc := range p.SampleTests {
			sb.WriteString(fmt.Sprintf("**Example %d:**\n\n", i+1))
			sb.WriteString("```\n")
			sb.WriteString("Input:  " + tc.Input + "\n")
			sb.WriteString("Output: " + tc.Output + "\n")
			sb.WriteString("```\n\n")
		}
	}

	// Hints
	if len(p.Hints) > 0 {
		sb.WriteString("## Hints\n\n")
		for i, hint := range p.Hints {
			sb.WriteString(fmt.Sprintf("<details>\n<summary>Hint %d</summary>\n\n%s\n\n</details>\n\n", i+1, htmlToMarkdown(hint)))
		}
	}

	// Notes section for personal use
	sb.WriteString("## Notes\n\n")
	sb.WriteString("<!-- Add your approach, complexity analysis, and reflections here -->\n\n")
	sb.WriteString("**Time complexity:** O(?)\n\n")
	sb.WriteString("**Space complexity:** O(?)\n")

	return sb.String()
}

// htmlToMarkdown does a best-effort conversion of LeetCode HTML to Markdown.
func htmlToMarkdown(h string) string {
	// Unescape HTML entities first
	s := html.UnescapeString(h)

	// Code blocks: <pre> content → fenced code blocks
	rePre := regexp.MustCompile(`(?s)<pre[^>]*>(.*?)</pre>`)
	s = rePre.ReplaceAllStringFunc(s, func(m string) string {
		inner := rePre.FindStringSubmatch(m)[1]
		inner = stripTags(inner)
		inner = strings.TrimSpace(html.UnescapeString(inner))
		return "```\n" + inner + "\n```"
	})

	// Inline code
	reCode := regexp.MustCompile(`<code>(.*?)</code>`)
	s = reCode.ReplaceAllString(s, "`$1`")

	// Bold
	reBold := regexp.MustCompile(`<strong>(.*?)</strong>`)
	s = reBold.ReplaceAllString(s, "**$1**")

	reBold2 := regexp.MustCompile(`<b>(.*?)</b>`)
	s = reBold2.ReplaceAllString(s, "**$1**")

	// Italic
	reEm := regexp.MustCompile(`<em>(.*?)</em>`)
	s = reEm.ReplaceAllString(s, "*$1*")

	// Lists
	reLi := regexp.MustCompile(`<li>(.*?)</li>`)
	s = reLi.ReplaceAllString(s, "- $1")

	// Paragraphs and line breaks → newlines
	reP := regexp.MustCompile(`</p>`)
	s = reP.ReplaceAllString(s, "\n\n")
	reBr := regexp.MustCompile(`<br\s*/?>`)
	s = reBr.ReplaceAllString(s, "\n")

	// Super/subscript
	reSup := regexp.MustCompile(`<sup>(.*?)</sup>`)
	s = reSup.ReplaceAllString(s, "^$1")
	reSub := regexp.MustCompile(`<sub>(.*?)</sub>`)
	s = reSub.ReplaceAllString(s, "_$1")

	// Strip remaining HTML tags
	s = stripTags(s)

	// Collapse excessive blank lines
	reBlank := regexp.MustCompile(`\n{3,}`)
	s = reBlank.ReplaceAllString(s, "\n\n")

	return strings.TrimSpace(s)
}

// stripTags removes all HTML tags.
func stripTags(s string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return re.ReplaceAllString(s, "")
}
