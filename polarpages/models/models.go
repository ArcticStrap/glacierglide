package models

import (
	"html/template"
)

type WebPage struct {
	Title   string
	Content template.HTML
	Theme   string
}

type PageEdit struct {
	EditId      int64  `json:"editId"`
	PageId      int64  `json:"pageId"`
	Date        string `json:"changeDate"`
	Time        string `json:"changeTime"`
	Editor      string `json:"editor"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

type SessionData struct {
	LoggedIn bool
	Username string
}

type StaticPage struct {
	Theme string
}

type WebModes struct {
	PageMode string
}
