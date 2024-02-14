package markdown

import "testing"

func TestMarkdownEngine(t *testing.T) {
  if len(Tokenize("# test")) != 1 {
    t.Fatal("Tokenizer did not detect header syntax")
  }
}
