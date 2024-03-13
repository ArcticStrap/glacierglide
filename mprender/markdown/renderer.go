package markdown

import (
	"fmt"
	"strconv"
)

func renderInline(p *Part) string {
	renderedLine := ""
	for _, v := range p.Children {
		switch vEle := v.(type) {
		case Code:
			renderedLine += fmt.Sprintf("<code>%s</code>", vEle.Value)
			break
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
		case Email:
			renderedLine += fmt.Sprintf("<a href=\"mailto:%s\">%s</a>", vEle.Path, vEle.Value)
		case Link:
			renderedLine += fmt.Sprintf("<a href=\"%s\" alt=\"%s\">%s</a>", vEle.Path, vEle.Alt,vEle.Value)
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
		case HorizontalRule:
			tTag = "hr"
			htmlOut += "<hr>\n"
		case BlockQuote:
			tTag = "blockquote"
			vValue = vPart.Value
			break
		case List:
			if vPart.Ordered {
				tTag = "ol"
			} else {
				tTag = "ul"
			}
			for i, ele := range vPart.Children {
				if i == 0 {
					vValue += "\n"
				}
				switch li := ele.(type) {
				case ListItem:
					vValue += fmt.Sprintf("<li>%s</li>\n", li.Value)
				}
			}
		case Paragraph:
			tTag = "p"
			vValue = renderInline(&vPart.Part)
			break
		default:
			tTag = "p"
			vValue = "ERROR: COULD NOT DETERMINE TYPE OF BLOCK"
			break
		}
		if tTag != "hr" {
			htmlOut += fmt.Sprintf("<%s>%s</%s>\n", tTag, vValue, tTag)
		}
	}

	return htmlOut
}
