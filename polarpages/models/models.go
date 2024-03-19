package models

import (
	"html/template"
)

type WebPage struct {
	Title   string
	Content template.HTML
	Theme   string
}

type SessionData struct {
	LoggedIn bool
	Username string
}

type StaticPage struct {
	Theme string
}
