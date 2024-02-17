package markdown

import "regexp"

// Struct setup

type Block interface {
}

type Part struct {
	Value string

	// Metadata
	Children []Block
}

type Header struct {
	Part

	Level int
}

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

func ParseInline(line []byte) []Block {
	var pChildren []Block

	bStart := 0
	for i := 0; i <= len(line); i++ {
		if i < len(line) && line[i] == '*' {
			// Check for single, double, or triple emphasis
			if i+1 < len(line) && line[i+1] == '*' {
				// Double emphasis
				if i+2 < len(line) && line[i+2] == '*' {
					// Bold Italic
					offset := i + 3
					for offset < len(line) {
						if offset < len(line) && line[offset] == '*' {
							if offset+1 < len(line) && line[offset+1] == '*' {
								if offset+2 < len(line) && line[offset+2] == '*' {
									break
								}
							}
						}
					}
					if offset != 0 && offset < len(line) {
						pChildren = append(pChildren, Bold{Part{Value: "", Children: []Block{Italic{Part{Value: string(line[i+3 : offset])}}}}})
					}
					i = offset + 3
					bStart = i
				} else {
					// Bold
					offset := i + 2
					for offset < len(line) {
						if offset < len(line) && line[offset] == '*' {
							if offset+1 < len(line) && line[offset+1] == '*' {
								break
							}
						}
						offset++
					}
					if offset != 0 && offset < len(line) {
						pChildren = append(pChildren, Bold{Part{Value: string(line[i+2 : offset])}})
					}
					i = offset + 2
					bStart = i
				}
			} else {
				// Italic
				offset := i + 1
				for offset < len(line) {
					if offset < len(line) && line[offset] == '*' {
						break
					}
					offset++
				}
				if offset != 0 && offset < len(line) {
					pChildren = append(pChildren, Italic{Part{Value: string(line[i+1 : offset])}})
				}
				i = offset + 1
				bStart = i
			}
		} else {
			// Add remaining text
			pChildren = append(pChildren, PlainText{Part{Value: string(line[bStart])}})
		}
	}
	return pChildren
}

/*func ParseBlockScope(content []byte) Block {

}*/

func Tokenize(content []byte) []Block {
	if len(content) == 0 {
		return nil
	}
	var blocks []Block

	headerRegex := regexp.MustCompile(`^(#+)\s+(.*)$`)

	bStart := 0
	parseMode := false

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
				bStart = i
				parseMode = false
				continue
			} else {
				blocks = append(blocks, Paragraph{Part: Part{Value: string(substr), Children: ParseInline(substr)}})
			}
		}

		if len(substr) != 0 && substr[len(substr)-1] == '#' && !parseMode {
			bStart = i - 1
			parseMode = true
		}
	}

	return blocks
}
