package leetcode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const graphqlURL = "https://leetcode.com/graphql"

// Problem holds all fetched data for a LeetCode problem.
type Problem struct {
	Slug        string
	FrontendID  string
	Title       string
	Difficulty  string
	Content     string // HTML
	Topics      []string
	CodeSnippet string // Go snippet
	SampleTests []SampleTest
	Hints       []string
}

type SampleTest struct {
	Input  string
	Output string
	Args   []string
}

type graphqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

type graphqlResponse struct {
	Data struct {
		Question struct {
			QuestionFrontendId string   `json:"questionFrontendId"`
			Title              string   `json:"title"`
			TitleSlug          string   `json:"titleSlug"`
			Difficulty         string   `json:"difficulty"`
			Content            string   `json:"content"`
			Hints              []string `json:"hints"`
			TopicTags          []struct {
				Name string `json:"name"`
			} `json:"topicTags"`
			CodeSnippets []struct {
				LangSlug string `json:"langSlug"`
				Code     string `json:"code"`
			} `json:"codeSnippets"`
			ExampleTestcaseList []string `json:"exampleTestcaseList"`
			MetaData            string   `json:"metaData"`
		} `json:"question"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

const problemQuery = `
query getProblem($titleSlug: String!) {
  question(titleSlug: $titleSlug) {
    questionFrontendId
    title
    titleSlug
    difficulty
    content
    hints
    topicTags { name }
    codeSnippets { langSlug code }
    exampleTestcaseList
    metaData
  }
}`

// FetchProblem fetches problem data from LeetCode by slug.
// Slug is the last segment of the problem URL, e.g. "two-sum".
func FetchProblem(slug string) (*Problem, error) {
	slug = NormalizeSlug(slug)

	body, err := json.Marshal(graphqlRequest{
		Query:     problemQuery,
		Variables: map[string]any{"titleSlug": slug},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", graphqlURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://leetcode.com/problems/"+slug+"/")
	req.Header.Set("User-Agent", "Mozilla/5.0 go-algo-workbench/1.0")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var gResp graphqlResponse
	if err := json.Unmarshal(data, &gResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	if len(gResp.Errors) > 0 {
		return nil, fmt.Errorf("LeetCode API error: %s", gResp.Errors[0].Message)
	}

	q := gResp.Data.Question
	if q.TitleSlug == "" {
		return nil, fmt.Errorf("problem %q not found", slug)
	}

	p := &Problem{
		Slug:       q.TitleSlug,
		FrontendID: q.QuestionFrontendId,
		Title:      q.Title,
		Difficulty: q.Difficulty,
		Content:    q.Content,
		Hints:      q.Hints,
	}

	for _, tag := range q.TopicTags {
		p.Topics = append(p.Topics, tag.Name)
	}

	for _, snippet := range q.CodeSnippets {
		if snippet.LangSlug == "golang" {
			p.CodeSnippet = snippet.Code
			break
		}
	}

	// Parse example test cases using metadata and HTML content for expected outputs.
	p.SampleTests = parseExampleTests(q.ExampleTestcaseList, q.MetaData, q.Content)

	return p, nil
}

// NormalizeSlug extracts slug from a full URL or returns as-is.
func NormalizeSlug(input string) string {
	// Handle full URLs like https://leetcode.com/problems/two-sum/
	if strings.Contains(input, "leetcode.com/problems/") {
		parts := strings.Split(input, "/problems/")
		if len(parts) > 1 {
			slug := strings.Split(parts[1], "/")[0]
			slug = strings.TrimSpace(slug)
			return slug
		}
	}
	return strings.TrimSpace(input)
}
