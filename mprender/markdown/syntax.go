package markdown

func ParseCode(text []byte, start int) ([]Chunk, int) {
	if len(text) < 2 {
		return []Chunk{Paragraph{Part{Value: string(text)}}}, 0
	}
	pChildren := []Chunk{}
	var pBytes []byte
	i := start

	// Append inactive plain text
	pChildren = append(pChildren, PlainText{Part{Value: string(text[0:i])}})

	for ; i < len(text); i++ {
		for i < len(text) && text[i] == ' ' {
			i++
		}
		if i >= len(text) {
			break
		}
		if text[i] == '`' {
			// Process content
			i++
			for i < len(text) && text[i] != '`' {
				pBytes = append(pBytes, text[i])
				i++
			}

			if i < len(text) {
				i++
				break
			}
		}
	}

	pChildren = append(pChildren, Code{Part{Value: string(pBytes)}})

	return pChildren, i - start
}

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

func ParseHeader(text []byte) Chunk {
	if len(text) < 3 {
		return Paragraph{Part{Value: string(text)}}
	}

	i := 0
	level := 0

	for ; i < len(text); i++ {
		for i < len(text) && text[i] == '#' {
			i++
			level++
		}
		if text[i] != ' ' || level > 6 {
			return Paragraph{Part{Value: string(text)}}
		}
		return Header{Part: Part{Value: string(text[i+1 : len(text)-1])}, Level: level}
	}

	return Paragraph{Part{Value: string(text)}}
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

func ParseLink(text []byte, start int) ([]Chunk, int) {
	pChildren := []Chunk{}

	i := start

	prefix := byte(' ')
	var title string
	var path string
	var alt string

	// Append inactive plain text
	if start > 0 {
		pChildren = append(pChildren, PlainText{Part{Value: string(text[0:i])}})
	}

	if start-1 > 0 {
		prefix = text[start-1]
	}

	// Skip '[' char
	i++

	tStart := i
	for i < len(text) && text[i] != ']' {
		i++
	}

	if i-1 > tStart {
		title = string(text[tStart:i])
	}

	if i+1 < len(text) && text[i+1] != '(' {
		pChildren = append(pChildren, PlainText{Part{Value: string(text[start:])}})
		return pChildren, i - start
	}

	i += 2

	pStart := i
	for i < len(text) && text[i] != ' ' && text[i] != ')' {
		i++
	}

	if i-1 > pStart {
		path = string(text[pStart:i])
	}

	// Get alt if two args separated by space
	if text[i] == ' ' {
		i++
		aStart := i

		for i < len(text) && text[i] != ')' {
			i++
		}
		if i-1 > aStart {
			alt = string(text[aStart:i])
		}
	}

	switch prefix {
	case '!':
		pChildren = append(pChildren, Image{Part: Part{Value: title}, Path: path, Alt: alt})
		break
	case ' ':
		pChildren = append(pChildren, Link{Part: Part{Value: title}, Path: path, Alt: alt})
		break
	// Append as plain text if invalid prefix
	default:
		pChildren = append(pChildren, PlainText{Part{Value: string(text[start:])}})
		break
	}

	i++

	return pChildren, i - start
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
