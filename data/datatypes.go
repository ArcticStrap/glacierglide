package data

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AccountReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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

func (u *User) ValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw)) == nil
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
