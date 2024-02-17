package markdown

import (
	"fmt"
	"strconv"
)

func renderInline(p *Part) string {
  renderedLine := ""
  for _,v := range p.Children {
    switch vEle := v.(type) {
    case Italic:
      renderedLine += fmt.Sprintf("<i>%s</i>",vEle.Value)
    case PlainText:
      renderedLine += vEle.Value
    default:
      renderedLine += "ERROR: COULD NOT DETERMINE TYPE OF ELEMENT"
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
      println("CHILDREN: ",len(vPart.Children))
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
