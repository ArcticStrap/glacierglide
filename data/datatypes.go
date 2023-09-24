package data

import "time"

type Page struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type User struct {
	UserId       int64
	Username     string
	Password     string
	CreationDate time.Time
}

type PageDiff struct {
	DiffId      int64
	PageId      int64
	Date        time.Time
	Time        time.Time
	UserId      int64
	Anon        bool
	Description string
	Content     string
}
