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
		// Header cases
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
		// Other cases
		{
			input:    "Horizontal rule below\n---",
			expected: "<p>Horizontal rule below</p>\n<hr>\n",
		},
		{
			input:    "Here is some `code` text.",
			expected: "<p>Here is some <code>code</code> text.</p>\n",
		},
		// Blockquote cases
		{
			input:    "> Single-line blockquote",
			expected: "<blockquote>Single-line blockquote</blockquote>\n",
		},
		{
			input:    "> First line of a multi-line\n> blockquote\n> Last line of a multi-line blockquote",
			expected: "<blockquote>First line of a multi-line\nblockquote\nLast line of a multi-line blockquote</blockquote>\n",
		},
		{
			input:    "> Blockquote with **bold** and *italic* text.",
			expected: "<blockquote>Blockquote with <strong>bold</strong> and <em>italic</em> text.</blockquote>\n",
		},
		{
			input:    "> Nested blockquote\n>> Nested within the first one",
			expected: "<blockquote>Nested blockquote\n<blockquote>Nested within the first one</blockquote>\n</blockquote>\n",
		},
		{
			input:    "> Blockquote with\n> multiple\n> lines\n> of text",
			expected: "<blockquote>Blockquote with\nmultiple\nlines\nof text</blockquote>\n",
		},
		// List cases
		{
			input:    "This is an ordered list:\n1. one\n2. two\n3. three",
			expected: "<p>This is an ordered list:</p>\n<ol>\n<li>one</li>\n<li>two</li>\n<li>three</li>\n</ol>\n",
		},
		{
			input:    "This is an unordered list:\n- 1\n- 2\n- 3",
			expected: "<p>This is an unordered list:</p>\n<ul>\n<li>1</li>\n<li>2</li>\n<li>3</li>\n</ul>\n",
		},
		// Emphasis cases
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
		// Link cases
		{
			input:    "This is a link to [wikipedia](https://wikipedia.org wikipedia).",
			expected: "<p>This is a link to <a href=\"https://wikipedia.org\" alt=\"wikipedia\">wikipedia</a>.</p>\n",
		},
		{
			input:    "This is a quick link: <https://wikipedia.org>.",
			expected: "<p>This is a quick link: <a href=\"https://wikipedia.org\" alt=\"\">https://wikipedia.org</a>.</p>\n",
		},
		{
			input:    "This is not a link: <>.",
			expected: "<p>This is not a link: <>.</p>\n",
		},
		{
			input:    "This is an email: <business@arcticstrap.net>.",
			expected: "<p>This is an email: <a href=\"mailto:business@arcticstrap.net\">business@arcticstrap.net</a>.</p>\n",
		},
		{
			input:    "Plain text",
			expected: "<p>Plain text</p>\n",
		},
		// Backslash cases
		{
			input:    "Backslashes: \\\\a",
			expected: "<p>Backslashes: \\a</p>\n",
		},
		{
			input:    "Show Backslashes: \\\\",
			expected: "<p>Show Backslashes: \\</p>\n",
		},
		{
			input:    "Escaped asterik: \\*",
			expected: "<p>Escaped asterik: *</p>\n",
		},
		{
			input:    "Escaped markup: \\**hello**",
			expected: "<p>Escaped markup: *<em>hello</em>*</p>\n",
		},
		{
			input:    "\\## Escaped header",
			expected: "<p>## Escaped header</p>\n",
		},
	}

	pass := true

	for _, tc := range testCases {
		actual := ToHTML(tc.input)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Logf("ERROR\nINPUT: %s\nEXPECTED: %s\nRESULT: %s\n", tc.input, tc.expected, actual)
			pass = false
		}
	}

	if pass {
		t.Log("TEXT PROPERLY RENDERED")
	} else {
		t.Errorf("FAIL")
	}
}
