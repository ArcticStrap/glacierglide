package markdown

import "regexp"

// Struct setup

type Block interface {
}

type Part struct {
	Value string
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

func Tokenize(content []byte) []Block {
	if len(content) == 0 {
		return nil
	}
	var blocks []Block

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
				blocks = append(blocks, Paragraph{Part: Part{Value: string(substr)}})
			}
		}

		if len(substr) != 0 && substr[len(substr)-1] == '#' && !parseMode {
			bStart = i - 1
			parseMode = true
		}
	}

	return blocks
}
