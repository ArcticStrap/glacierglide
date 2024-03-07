package models

import (
	"html/template"
)

type WebPage struct {
	Title   string
	Content template.HTML
	Theme   string
}

type StaticPage struct {
	Theme string
}
