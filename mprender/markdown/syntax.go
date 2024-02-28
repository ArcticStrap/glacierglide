package markdown

func ParseBlockQuote(text []byte) []Chunk {
	if len(text) < 2 {
		return nil
	}
	pChildren := []Chunk{}

	for i := 0; i < len(text); i++ {
		if (text[i] == '\n' && i-1 >= 0 && text[i-1] == '\n') || i == len(text)-1 {
			sCount := 1
			if text[1] == ' ' {
				sCount++
			}
			pChildren = append(pChildren, BlockQuote{Part{Value: string(text[sCount:i])}})
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
