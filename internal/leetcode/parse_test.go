package leetcode

import "testing"

func TestParseExampleTests_FillsArgsAndOutputs(t *testing.T) {
	inputs := []string{
		"[\"flower\",\"flow\",\"flight\"]",
		"[\"dog\",\"racecar\",\"car\"]",
	}
	rawMeta := `{"params":[{"name":"strs","type":"string[]"}],"return":{"type":"string"}}`
	content := `<pre><strong>Input:</strong> strs = [&quot;flower&quot;,&quot;flow&quot;,&quot;flight&quot;]
<strong>Output:</strong> &quot;fl&quot;
</pre>
<pre><strong>Input:</strong> strs = [&quot;dog&quot;,&quot;racecar&quot;,&quot;car&quot;]
<strong>Output:</strong> &quot;&quot;
</pre>`

	tests := parseExampleTests(inputs, rawMeta, content)
	if len(tests) != 2 {
		t.Fatalf("len(tests) = %d, want 2", len(tests))
	}
	if tests[0].Input != `strs = ["flower","flow","flight"]` {
		t.Fatalf("first input = %q", tests[0].Input)
	}
	if len(tests[0].Args) != 1 || tests[0].Args[0] != `["flower","flow","flight"]` {
		t.Fatalf("first args = %#v", tests[0].Args)
	}
	if tests[0].Output != `"fl"` {
		t.Fatalf("first output = %q", tests[0].Output)
	}
	if tests[1].Output != `""` {
		t.Fatalf("second output = %q", tests[1].Output)
	}
}

func TestParseExampleTests_ParsesArgsFromExampleList(t *testing.T) {
	inputs := []string{
		"[1,2,3,1]",
		"[1,2,3,4]",
	}
	rawMeta := `{"params":[{"name":"nums","type":"integer[]"}],"return":{"type":"boolean"}}`
	content := `<div class="example-block">
<p><strong>Input:</strong> <span class="example-io">nums = [1,2,3,1]</span></p>
<p><strong>Output:</strong> <span class="example-io">true</span></p>
</div>`

	tests := parseExampleTests(inputs, rawMeta, content)
	if len(tests) != 2 {
		t.Fatalf("len(tests) = %d, want 2", len(tests))
	}
	if len(tests[0].Args) != 1 || tests[0].Args[0] != "[1,2,3,1]" {
		t.Fatalf("first args = %#v", tests[0].Args)
	}
	if tests[0].Output != "true" {
		t.Fatalf("first output = %q, want %q", tests[0].Output, "true")
	}
}

func TestExtractOutputFromExampleBlock_MultilineOutput(t *testing.T) {
	block := "Input: x = 1\nOutput:\n[1, 2, 3]\nExplanation: sample"
	got := extractOutputFromExampleBlock(block)
	if got != "[1, 2, 3]" {
		t.Fatalf("output = %q, want %q", got, "[1, 2, 3]")
	}
}
