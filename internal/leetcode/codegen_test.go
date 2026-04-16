package leetcode

import (
	"strings"
	"testing"
)

func TestToGoLiteral(t *testing.T) {
	tests := []struct {
		name   string
		raw    string
		goType string
		want   string
		ok     bool
	}{
		{name: "int", raw: "42", goType: "int", want: "42", ok: true},
		{name: "string quoted", raw: `"abc"`, goType: "string", want: `"abc"`, ok: true},
		{name: "string bare", raw: "abc", goType: "string", want: `"abc"`, ok: true},
		{name: "slice", raw: "[1,2,3]", goType: "[]int", want: "[]int{1,2,3}", ok: true},
		{name: "null slice", raw: "null", goType: "[]int", want: "nil", ok: true},
		{name: "null int unsupported", raw: "null", goType: "int", want: "", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := toGoLiteral(tt.raw, tt.goType)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("got = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTypedTest_UsesParsedArgs(t *testing.T) {
	g := &GoCodegen{
		Problem: &Problem{
			SampleTests: []SampleTest{
				{
					Input:  "nums = [1,2,3,1]",
					Output: "true",
					Args:   []string{"[1,2,3,1]"},
				},
			},
		},
	}

	rendered := g.typedTest("containsduplicate", "containsDuplicate", []param{
		{name: "nums", goType: "[]int"},
	}, "bool")

	if !strings.Contains(rendered, "nums: []int{1,2,3,1},") {
		t.Fatalf("typed test did not include parsed args:\\n%s", rendered)
	}
	if !strings.Contains(rendered, "want: true,") {
		t.Fatalf("typed test did not include parsed output:\\n%s", rendered)
	}
}

func TestJSONArrayToGoBraces(t *testing.T) {
	got, ok := jsonArrayToGoBraces(`[[1,2],["[ok]"]]`)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	want := `{{1,2},{"[ok]"}}`
	if got != want {
		t.Fatalf("got = %q, want %q", got, want)
	}
}
