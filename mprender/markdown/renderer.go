package markdown

import "fmt"

func ToHTML(content string) string {
	tokens := Tokenize(content)
	htmlOut := ""

	for _, v := range tokens {
		tTag := "p"
		switch v.Type {
		case Header1:
			tTag = "h1"
			break
		case Header2:
			tTag = "h2"
			break
		case Header3:
			tTag = "h3"
			break
		case Header4:
			tTag = "h4"
			break
		case Header5:
			tTag = "h5"
			break
		case Header6:
			tTag = "h6"
			break
		case Bold:
			tTag = "strong"
		case Italic:
			tTag = "em"
		case Paragraph:
			tTag = "p"
			break
		default:
			tTag = "p"
			break
		}
		htmlOut += fmt.Sprintf("<%s>%s</%s>\n", tTag, v.Value, tTag)
	}

	return htmlOut
}
