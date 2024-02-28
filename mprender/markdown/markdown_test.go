package markdown

import (
	"reflect"
	"testing"
)

func TestMarkdownEngine(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "# Header 1",
			expected: "<h1>Header 1</h1>\n",
		},
		{
			input:    "## Header 2",
			expected: "<h2>Header 2</h2>\n",
		},
		{
			input:    "### Header 3",
			expected: "<h3>Header 3</h3>\n",
		},
		{
			input:    "#### Header 4",
			expected: "<h4>Header 4</h4>\n",
		},
		{
			input:    "##### Header 5",
			expected: "<h5>Header 5</h5>\n",
		},
		{
			input:    "###### Header 6",
			expected: "<h6>Header 6</h6>\n",
		},
    {
      input: "> This is a quote block.",
      expected: "<blockquote>This is a quote block.</blockquote>\n",
    },
		{
			input:    "This is *italic* text.",
			expected: "<p>This is <em>italic</em> text.</p>\n",
		},
		{
			input:    "This is **bold** text.",
			expected: "<p>This is <strong>bold</strong> text.</p>\n",
		},
		{
			input:    "This is ***bold* and *italic*** text.",
			expected: "<p>This is <strong><em>bold</em> and <em>italic</em></strong> text.</p>\n",
		},
		{
			input:    "Plain text",
			expected: "<p>Plain text</p>\n",
		},
	}

	pass := true

	for _, tc := range testCases {
		actual := ToHTML(tc.input)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Logf("INPUT: %s\nEXPECTED: %s\nRESULT: %s\n", tc.input, tc.expected, actual)
			pass = false
		}
	}

	if pass {
		t.Log("TEXT PROPERLY RENDERED")
	} else {
		t.Errorf("FAIL")
	}
}
