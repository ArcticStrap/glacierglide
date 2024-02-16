package markdown

import "testing"

func TestMarkdownEngine(t *testing.T) {
	test_case := "# test 1\n## test 2\n\n**bold**\n\n*italic*"

	tokens := Tokenize([]byte(test_case))
	if len(tokens) != 6 {
		t.Fatalf("Tokenizer did not detect header syntax: %d != 2", len(tokens))
	} else {
		for i, v := range tokens {
			t.Logf("TOKEN %d: TYPE: %d, VALUE: %s", i, v.Type, v.Value)
		}
	}

  t.Log(ToHTML(test_case))
}
