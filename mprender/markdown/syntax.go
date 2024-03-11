package markdown

import "regexp"

func ParseBlockQuote(text []byte) ([]Chunk, int) {
	if len(text) < 2 {
		return []Chunk{Paragraph{Part{Value: string(text)}}}, 0
	}
	pChildren := []Chunk{}
	var pBytes []byte
	i := 0

	for ; i < len(text); i++ {
		for i < len(text) && text[i] == ' ' {
			i++
		}
		if i >= len(text) {
			break
		}
		if text[i] == '>' {
			// Add new line if not first blockquote
			if i != 0 {
				pBytes = append(pBytes, '\n')
			}

			// Skip > and spaces
			for i < len(text) && (text[i] == '>' || text[i] == ' ') {
				i++
			}
			// Process content
			for i < len(text) && text[i] != '\n' {
				pBytes = append(pBytes, text[i])
				i++
			}

			if i < len(text) {
				if text[i] == '\n' && i+1 < len(text) && text[i+1] == '\n' {
					i++
					break
				}
			}
		}
	}

	pChildren = append(pChildren, BlockQuote{Part{Value: string(pBytes)}})

	return pChildren, i
}

func ParseHeader(text []byte) []Chunk {
	if len(text) < 3 {
		return []Chunk{Paragraph{Part{Value: string(text)}}}
	}
	pChildren := []Chunk{}

	headerRegex := regexp.MustCompile(`^(#+)\s+(.*)$`)

	if headerMatch := headerRegex.FindStringSubmatch(string(text[:len(text)-1])); len(headerMatch) == 3 {
		// Check header level
		switch len(headerMatch[1]) {
		case 1:
			pChildren = append(pChildren, Header{Part: Part{Value: headerMatch[2]}, Level: 1})
			break
		case 2:
			pChildren = append(pChildren, Header{Part: Part{Value: headerMatch[2]}, Level: 2})
			break
		case 3:
			pChildren = append(pChildren, Header{Part: Part{Value: headerMatch[2]}, Level: 3})
			break
		case 4:
			pChildren = append(pChildren, Header{Part: Part{Value: headerMatch[2]}, Level: 4})
			break
		case 5:
			pChildren = append(pChildren, Header{Part: Part{Value: headerMatch[2]}, Level: 5})
			break
		case 6:
			pChildren = append(pChildren, Header{Part: Part{Value: headerMatch[2]}, Level: 6})
			break
		default:
			pChildren = append(pChildren, Paragraph{Part: Part{Value: headerMatch[2]}})
			break
		}
	}
	return pChildren
}

func ParseEmph(text []byte, start int, eChar byte) ([]Chunk, int) {
	pChildren := []Chunk{}
	i := start
	for ; i < len(text); i++ {
		// Count consecutive asteriks
		aCount := 1
		for i+aCount < len(text) && text[i+aCount] == eChar {
			aCount++
		}
		if aCount > 3 {
			aCount = 3
		}

		switch aCount {
		case 1:
			offset := i + 1
			for offset < len(text) && text[offset] != eChar {
				offset++
			}
			// Append inactive plain text
			pChildren = append(pChildren, PlainText{Part{Value: string(text[0:i])}})

			if offset != 0 && offset < len(text) && text[offset] == eChar {
				pChildren = append(pChildren, Italic{Part{Value: string(text[i+1 : offset])}})
			} else {
				pChildren = append(pChildren, PlainText{Part{Value: string(text[i+1 : offset])}})
			}
			i = offset + 1
			break
		case 2:
			offset := i + 2
			for offset < len(text) {
				if text[offset] == eChar {
					if offset+1 < len(text) && text[offset+1] == eChar {
						break
					}
				}
				offset++
			}
			// Append inactive plain text
			pChildren = append(pChildren, PlainText{Part{Value: string(text[0:i])}})
			if offset != 0 && offset < len(text) && text[offset] == eChar {
				pChildren = append(pChildren, Bold{Part{Value: string(text[i+2 : offset])}})
			} else {
				pChildren = append(pChildren, PlainText{Part{Value: string(text[i+2 : offset])}})
			}
			i = offset + 2
			break
		case 3:
			offset := i + 3
			for offset < len(text) {
				if text[offset] == eChar {
					if offset+1 < len(text) && text[offset+1] == eChar {
						if offset+2 < len(text) && text[offset+2] == eChar {
							break
						}
					}
				}
				offset++
			}
			// Append inactive plain text
			pChildren = append(pChildren, PlainText{Part{Value: string(text[0:i])}})

			if offset != 0 && offset < len(text) && text[offset] == eChar {
				if offset+2 < len(text) && text[offset+1] == eChar && text[offset+2] == eChar {
					pChildren = append(pChildren, Bold{Part{Value: "", Children: ParseInline(text[i+2 : offset+1])}})
				} else if offset+1 < len(text) && text[offset+1] == eChar {
					pChildren = append(pChildren, Italic{Part{Value: string(text[i+3 : offset])}})
				} else {
					pChildren = append(pChildren, Bold{Part{Value: string(text[i+3 : offset])}})
				}
			} else {
				pChildren = append(pChildren, PlainText{Part{Value: string(text[i+3 : offset])}})
			}

			i = offset + 3
			break
		}
		if len(pChildren) != 0 {
			break
		}
	}
	return pChildren, i - start
}

func ParseOList(text []byte) (Chunk, int) {
	pList := List{Part: Part{Children: []Chunk{}}, Ordered: true}
	i := 0
	for ; i < len(text); i++ {
		for i < len(text) && text[i] == ' ' {
			i++
		}
		if i >= len(text) {
			break
		}
		for i < len(text) && text[i] >= '0' && text[i] <= '9' {
			i++
		}
		if i >= len(text) {
			break
		}
		if text[i] == '.' {
			for i < len(text) && (text[i] == '.' || text[i] == ' ') {
				i++
			}
			start := i

			for i < len(text) && text[i] != '\n' {
				i++
			}

			pList.Children = append(pList.Children, ListItem{Part{Value: string(text[start:i])}})
			if i < len(text) {
				if text[i] == '\n' && i+1 < len(text) && text[i+1] == '\n' {
          i++ 
          break
				}
			}
    }
	}

	return pList, i
}

func ParseUList(text []byte) (Chunk, int) {
	if len(text) < 2 {
		return Paragraph{Part{Value: string(text)}}, 0
	}

	pList := List{Part: Part{Children: []Chunk{}}, Ordered: false}
	i := 0

	for ; i < len(text); i++ {
		for i < len(text) && text[i] == ' ' {
			i++
		}
		if i >= len(text) {
			break
		}
		if text[i] == '-' {
			for i < len(text) && (text[i] == '-' || text[i] == ' ') {
				i++
			}
			start := i

			for i < len(text) && text[i] != '\n' {
				i++
			}

			pList.Children = append(pList.Children, ListItem{Part{Value: string(text[start:i])}})
			if i < len(text) {
				if text[i] == '\n' && i+1 < len(text) && text[i+1] != '-' {
					i++
					break
				}
			}
		}
	}

	return pList, i
}

func ParseQuickLink(text []byte, start int) ([]Chunk, int) {
	pChildren := []Chunk{}

	i := start

	// Email flag
	email := false

	// Append inactive plain text
	if start > 0 {
		pChildren = append(pChildren, PlainText{Part{Value: string(text[0:i])}})
	}

	for ; i < len(text); i++ {
		if text[i] == '>' {
			break
		} else if text[i] == '@' {
			email = true
		}
	}

	if i > len(text) || text[i] != '>' || len(text[start:i]) < 3 {
		pChildren = append(pChildren, PlainText{Part{Value: string(text[start : i+1])}})
		return pChildren, (i - start) + 1
	}

	if start+1 < i {
		if email {
			pChildren = append(pChildren, Email{Part: Part{Value: string(text[start+1 : i])}, Path: string(text[start+1 : i])})
		} else {
			pChildren = append(pChildren, Link{Part: Part{Value: string(text[start+1 : i])}, Path: string(text[start+1 : i])})
		}
	}

	return pChildren, (i - start) + 1
}
