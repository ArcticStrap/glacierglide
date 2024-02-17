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
	for i := 0; i < len(line); i++ {
		if line[i] == '*' {
			// Count consecutive asteriks
			aCount := 1
			for i+aCount < len(line) && line[i+aCount] == '*' {
				aCount++
			}
			if aCount > 3 {
				aCount = 3
			}

			switch aCount {
			case 1:
				offset := i + 1
				for offset < len(line) && line[offset] != '*' {
					offset++
				}
				// Append inactive plain text
				pChildren = append(pChildren, PlainText{Part{Value: string(line[bStart:i])}})
				if offset != 0 && offset < len(line) && line[offset] == '*' {
					pChildren = append(pChildren, Italic{Part{Value: string(line[i+1 : offset])}})
				}
				i = offset + 1
			case 2:
				offset := i + 2
				for offset < len(line) && line[offset] != '*' {
					offset++
				}
				// Append inactive plain text
				pChildren = append(pChildren, PlainText{Part{Value: string(line[bStart:i])}})
				if offset != 0 && offset < len(line) && line[offset] == '*' {
					pChildren = append(pChildren, Bold{Part{Value: string(line[i+2 : offset])}})
				}
				i = offset + 2
			case 3:
				offset := i + 3
				for offset < len(line) && line[offset] != '*' {
					offset++
				}
				// Append inactive plain text
				pChildren = append(pChildren, PlainText{Part{Value: string(line[bStart:i])}})
				if offset != 0 && offset < len(line) && line[offset] == '*' {
					pChildren = append(pChildren, Bold{Part{Value: "", Children: []Block{Italic{Part{Value: string(line[i+2 : offset])}}}}})
				}
				i = offset + 3
			}
			bStart = i
		} else {
			i++
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
