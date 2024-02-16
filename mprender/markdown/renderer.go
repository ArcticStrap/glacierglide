package markdown

import (
	"fmt"
	"strconv"
)

func ToHTML(content string) string {
	blocks := Tokenize([]byte(content))
	htmlOut := ""

	for _, v := range blocks {
		tTag := "p"
    var vValue string
    switch vPart := v.(type) {
		case Header:
			tTag = "h" + strconv.Itoa(vPart.Level)
      vValue = vPart.Value
			break
		case Bold:
			tTag = "strong"
      vValue = vPart.Value
			break
		case Italic:
			tTag = "em"
      vValue = vPart.Value
			break
		case Paragraph:
			tTag = "p"
      vValue = vPart.Value
			break
		default:
			tTag = "p"
      vValue = "ERROR: COULD NOT DETERMINE TYPE of Part"
			break
		}
		htmlOut += fmt.Sprintf("<%s>%s</%s>\n", tTag, vValue, tTag)
	}

	return htmlOut
}
