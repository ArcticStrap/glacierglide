package markdown

import "regexp"

// Struct setup

type Chunk interface {
}

type Part struct {
	Value string

	// Metadata
	Children []Chunk
}

type Header struct {
	Part

	Level int
}

type BlockQuote struct {
	Part
}

// Emphasis

type Bold struct {
	Part
}

type Italic struct {
	Part
}

type Paragraph struct {
	Part
}

type PlainText struct {
	Part
}

func ParseInline(line []byte) []Chunk {
	var pChildren []Chunk

	bStart := 0
	for i := 0; i < len(line); i++ {
		if line[i] == '*' || line[i] == '_' {
			parts, jump := ParseEmph(line[bStart:], i-bStart, line[i])
			if jump == 0 {
				continue
			}
			pChildren = append(pChildren, parts...)
			i += jump
			bStart = i
		}
	}

	// Append remaining text if any
	charCount := len(line)
	if bStart < charCount {
		if line[charCount-1] == ' ' || line[charCount-1] == '\n' {
			charCount--
		}
		pChildren = append(pChildren, PlainText{Part{Value: string(line[bStart:charCount])}})
	}

	return pChildren
}

func Tokenize(content []byte) []Chunk {
	if len(content) == 0 {
		return nil
	}
	var blocks []Chunk

	headerRegex := regexp.MustCompile(`^(#+)\s+(.*)$`)

	bStart := 0

	for i := 0; i <= len(content); i++ {
		substr := content[bStart:i]

		if len(substr) != 0 && (substr[len(substr)-1] == '\n' || substr[len(substr)-1] == '\r' || i == len(content)) {
			// Make it so the acutal last character can render in the parser
			if i == len(content) {
				substr = append(substr, ' ')
			}
			// Find header
			if headerMatch := headerRegex.FindStringSubmatch(string(substr[:len(substr)-1])); len(headerMatch) == 3 {
				// Check header level
				switch len(headerMatch[1]) {
				case 1:
					blocks = append(blocks, Header{Part: Part{Value: headerMatch[2]}, Level: 1})
					break
				case 2:
					blocks = append(blocks, Header{Part: Part{Value: headerMatch[2]}, Level: 2})
					break
				case 3:
					blocks = append(blocks, Header{Part: Part{Value: headerMatch[2]}, Level: 3})
					break
				case 4:
					blocks = append(blocks, Header{Part: Part{Value: headerMatch[2]}, Level: 4})
					break
				case 5:
					blocks = append(blocks, Header{Part: Part{Value: headerMatch[2]}, Level: 5})
					break
				case 6:
					blocks = append(blocks, Header{Part: Part{Value: headerMatch[2]}, Level: 6})
					break
				default:
					blocks = append(blocks, Paragraph{Part: Part{Value: headerMatch[2]}})
					break
				}
			} else if substr[0] == '>' {
				parts := ParseBlockQuote(substr[bStart:])
				blocks = append(blocks, parts...)
			} else {
				blocks = append(blocks, Paragraph{Part: Part{Value: string(substr), Children: ParseInline(substr)}})
			}
			bStart = i
		}
	}

	return blocks
}
