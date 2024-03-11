package markdown

// Credit: https://github.com/gomarkdown/markdown/blob/master/parser/block.go#L1772
func SkipChar(data []byte, start int, char byte) int {
	n := len(data)
	for start < n && data[start] == char {
		start++
	}
	return start
}

func SkipCharN(data []byte, start int, char byte, max int) int {
	n := len(data)
	for start < n && max > 0 && data[start] == char {
		start++
		max--
	}
	return start
}
