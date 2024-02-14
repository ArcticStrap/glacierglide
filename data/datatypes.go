package data

import (
	"encoding/json"
	"time"

	"github.com/ChaosIsFramecode/horinezumi/mprender"
	"golang.org/x/crypto/bcrypt"
)

type AccountReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Page struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	MPType  mprender.MPType `json:"mpType"`
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
	DiffId      int64     `pgx:"diff_id"     json:"diffId"`
	PageId      int64     `pgx:"page_id"     json:"pageId"`
	Date        time.Time `pgx:"change_date" json:"changeDate"`
	Time        time.Time `pgx:"change_time" json:"changeTime"`
	UserId      string    `pgx:"editor"     json:"editor"`
	Description string    `pgx:"description" json:"description"`
	Content     string    `pgx:"content"     json:"content"`
}

func (pd PageDiff) MarshalJSON() ([]byte, error) {
	type Alias PageDiff
	return json.Marshal(&struct {
		Date string `json:"changeDate"`
		Time string `json:"changeTime"`
		*Alias
	}{
		Date:  pd.Date.Format("2006-01-02"), // format date as "YYYY-MM-DD"
		Time:  pd.Time.Format("15:04:05"),   // format time as "HH:MM:SS"
		Alias: (*Alias)(&pd),
	})
}
