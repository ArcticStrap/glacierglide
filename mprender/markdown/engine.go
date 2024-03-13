package markdown

// Struct setup

type Chunk interface {
}

// Base

type Paragraph struct {
	Part
}

type PlainText struct {
	Part
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

type HorizontalRule struct {
}

type BlockQuote struct {
	Part
}

// Code
type Code struct {
	Part
}

type CodeBlock struct {
	Part

	Language string
}

// Emphasis

type Bold struct {
	Part
}

type Italic struct {
	Part
}

// Links
type Link struct {
	Part

	Path string
	Alt  string
}

type Email Link
type Image Link

// Lists
type List struct {
	Part

	Ordered bool
}

type ListItem struct {
	Part
}

func ParseInline(line []byte) []Chunk {
	var pChildren []Chunk

	bStart := 0
	for i := 0; i < len(line); i++ {
		if line[i] == '\\' && (i-1 <= 0 || line[i-1] != '\\') {
			// Ignore syntax if escape char is found
			if i+1 < len(line) {
				i++
				bStart++
			}
		} else if line[i] == '*' || line[i] == '_' {
			parts, jump := ParseEmph(line[bStart:], i-bStart, line[i])
			if jump == 0 {
				continue
			}
			pChildren = append(pChildren, parts...)
			i += jump
			bStart = i
		} else if line[i] == '<' {
			parts, jump := ParseQuickLink(line[bStart:], i-bStart)
			if jump == 0 {
				continue
			}

			pChildren = append(pChildren, parts...)
			i += jump
			bStart = i
		} else if line[i] == '`' {
			parts, jump := ParseCode(line[bStart:], i-bStart)
			if jump == 0 {
				continue
			}

			pChildren = append(pChildren, parts...)
			i += jump
			bStart = i
		} else if line[i] == '[' {
			parts, jump := ParseLink(line[bStart:], i-bStart)
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

	content = RmvCr(content)

	var blocks []Chunk

	bStart := 0

	for i := 0; i <= len(content); i++ {
		substr := content[bStart:i]

		if len(substr) != 0 && (substr[len(substr)-1] == '\n' || i == len(content)) {
			// Make it so the acutal last character can render in the parser
			if i == len(content) {
				substr = append(substr, ' ')
			}

			// Append as paragraph if escape char found
			if substr[0] == '\\' {
				blocks = append(blocks, Paragraph{Part: Part{Value: "", Children: ParseInline(substr[0:])}})
				continue
			}

			// Find blocks
			if substr[0] == '#' {
				nBlock := ParseHeader(substr)
				blocks = append(blocks, nBlock)
				bStart += len(substr)
				i = bStart
				continue
			} else if substr[0] == '>' {
				nBlocks, jump := ParseBlockQuote(content[bStart:])
				blocks = append(blocks, nBlocks...)
				bStart += jump
				i = bStart
				continue
			} else if substr[0] == '-' {
				// Check for potential horizontal rule
				if 1 < len(substr)-1 && substr[1] == '-' {
					if 2 < len(substr)-1 && substr[2] == '-' {
						jump := 2
						for jump < len(substr) && (substr[jump] == '-' || substr[jump] == ' ') {
							jump++
						}
						blocks = append(blocks, HorizontalRule{})
						bStart += jump
						i = bStart
						continue
					}
				}

				// Parse unordered list if otherwise
				nBlock, jump := ParseUList(content[bStart:])
				blocks = append(blocks, nBlock)
				bStart += jump
				i = bStart
				continue
			} else if substr[0] == '1' && 1 < len(substr)-1 && substr[1] == '.' {
				nBlock, jump := ParseOList(content[bStart:])
				blocks = append(blocks, nBlock)
				bStart += jump
				i = bStart
				continue
			} else {
				blocks = append(blocks, Paragraph{Part: Part{Value: "", Children: ParseInline(substr)}})
			}
			bStart = i
		}
	}

	return blocks
}

func RmvCr(str []byte) []byte {
	wi := 0
	strlen := len(str)
	for i := 0; i < strlen; i++ {
		char := str[i]

		// Continue if not \r
		if char != 13 {
			str[wi] = char
			wi++
			continue
		}

		// Replace \r with \n
		str[wi] = 10
		wi++
		if i < strlen-1 && str[i+1] == 10 {
			// If CRLF, skip \n
			i++
		}
	}

	return str[:wi]
}
