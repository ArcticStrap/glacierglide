package markdown

import "testing"

func TestMarkdownEngine(t *testing.T) {
	test_case := "# test 1\n## test 2\n\n**bold**\n\n*italic*"

	tokens := Tokenize([]byte(test_case))
	if len(tokens) != 6 {
		t.Fatalf("Tokenizer did not detect header syntax: %d != 2", len(tokens))
	}

  t.Log(ToHTML(test_case))
}
