package markdown

import "regexp"

// Struct setup

type TokenType int

const (
	Header1 TokenType = iota
	Header2
	Header3
	Header4
	Header5
	Header6
	Bold
	Italic
	Paragraph
)

type Token struct {
	Type  TokenType
	Value string
}

func Tokenize(content string) []Token {
	var tokens []Token

	headerRegex := regexp.MustCompile(`^(#+)\s+(.*)$`)
	italicRegex := regexp.MustCompile(`\*(.*?)\*`)
	boldRegex := regexp.MustCompile(`\*\*(.*?)\*\*`)

	lines := regexp.MustCompile(`\r\n|\n|\r`).Split(content, -1)
	for _, line := range lines {
		if headerMatch := headerRegex.FindStringSubmatch(line); len(headerMatch) == 3 {
			headerLevel := len(headerMatch[1])
			switch headerLevel {
			case 1:
				tokens = append(tokens, Token{Type: Header1, Value: headerMatch[2]})
				break
			case 2:
				tokens = append(tokens, Token{Type: Header2, Value: headerMatch[2]})
				break
			case 3:
				tokens = append(tokens, Token{Type: Header3, Value: headerMatch[2]})
				break
			case 4:
				tokens = append(tokens, Token{Type: Header4, Value: headerMatch[2]})
				break
			case 5:
				tokens = append(tokens, Token{Type: Header5, Value: headerMatch[2]})
				break
			case 6:
				tokens = append(tokens, Token{Type: Header6, Value: headerMatch[2]})
				break
			default:
				tokens = append(tokens, Token{Type: Paragraph, Value: line})
				break
			}
		} else if boldMatch := boldRegex.FindStringSubmatch(line); len(boldMatch) == 2 {
			tokens = append(tokens, Token{Type: Bold, Value: boldMatch[1]})
		} else if italicMatch := italicRegex.FindStringSubmatch(line); len(italicMatch) == 2 {
			tokens = append(tokens, Token{Type: Italic, Value: italicMatch[1]})
		} else {
			// Assume its a paragraph
			tokens = append(tokens, Token{Type: Paragraph, Value: line})
		}

	}

	return tokens
}
