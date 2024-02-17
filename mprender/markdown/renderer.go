package markdown

import (
	"fmt"
	"strconv"
)

func renderInline(p *Part) string {
	renderedLine := ""
	for _, v := range p.Children {
		switch vEle := v.(type) {
		case Italic:
			renderedLine += fmt.Sprintf("<em>%s</em>", vEle.Value)
			break
		case Bold:
      boldText := vEle.Value
      if len(vEle.Children) != 0 {
        boldText = renderInline(&vEle.Part)
      }
			renderedLine += fmt.Sprintf("<strong>%s</strong>", boldText)
			break
		case PlainText:
			renderedLine += vEle.Value
			break
		default:
			renderedLine += "ERROR: COULD NOT DETERMINE TYPE OF ELEMENT"
			break
		}
	}

	return renderedLine
}

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
		case Paragraph:
			tTag = "p"
			vValue = renderInline(&vPart.Part)
			break
		default:
			tTag = "p"
			vValue = "ERROR: COULD NOT DETERMINE TYPE OF BLOCK"
			break
		}
		htmlOut += fmt.Sprintf("<%s>%s</%s>\n", tTag, vValue, tTag)
	}

	return htmlOut
}
