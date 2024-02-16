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

func Tokenize(content []byte) []Token {
	if len(content) == 0 {
		return nil
	}
	var tokens []Token

	headerRegex := regexp.MustCompile(`^(#+)\s+(.*)$`)
	//italicRegex := regexp.MustCompile(`\*(.*?)\*`)
	//boldRegex := regexp.MustCompile(`\*\*(.*?)\*\*`)

	bStart := 0
	parseMode := false

	for i := 0; i <= len(content); i++ {
		substr := content[bStart:i]

		if len(substr) != 0 && (substr[len(substr)-1] == '\n' || substr[len(substr)-1] == '\r' || i == len(content)) {
			if i == len(content) {
				substr = append(substr, ' ')
			}
			if headerMatch := headerRegex.FindStringSubmatch(string(substr[:len(substr)-1])); len(headerMatch) == 3 {
				// Check header level
				switch len(headerMatch[1]) {
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
					tokens = append(tokens, Token{Type: Paragraph, Value: headerMatch[2]})
					break
				}
				bStart = i
				parseMode = false
				continue
			}
		}

		if len(substr) != 0 && substr[len(substr)-1] == '#' && !parseMode {
			bStart = i - 1
			parseMode = true
		}
	}

	return tokens
}
