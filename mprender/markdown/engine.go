package markdown

import "regexp"

// Struct setup

type TokenType int

const (
	Header1 TokenType = iota
	Paragraph
)

type Token struct {
	Type  TokenType
	Value string
}

func Tokenize(content string) []Token {
	var tokens []Token

	headerRegex := regexp.MustCompile(`^#\s+(.*)$`)

	lines := regexp.MustCompile(`\r\n|\n|\r`).Split(content, -1)
	for _, line := range lines {
		if headerMatch := headerRegex.FindStringSubmatch(line); len(headerMatch) == 2 {
			tokens = append(tokens, Token{Type: Header1, Value: headerMatch[1]})
		}
	}

	return tokens
}
