package data

import (
	"encoding/json"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AccountReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Page struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Namespace int    `json:"namespace"`
	MPType    int    `json:"mpType"`
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

type PageEdit struct {
	EditId      int64     `pgx:"edit_id"     json:"editId"`
	PageId      int64     `pgx:"page_id"     json:"pageId"`
	Date        time.Time `pgx:"change_date" json:"changeDate"`
	Time        time.Time `pgx:"change_time" json:"changeTime"`
	UserId      string    `pgx:"editor"     json:"editor"`
	Description string    `pgx:"description" json:"description"`
	Content     string    `pgx:"content"     json:"content"`
}

func (pe PageEdit) MarshalJSON() ([]byte, error) {
	type Alias PageEdit
	return json.Marshal(&struct {
		Date string `json:"changeDate"`
		Time string `json:"changeTime"`
		*Alias
	}{
		Date:  pe.Date.Format("2006-01-02"), // format date as "YYYY-MM-DD"
		Time:  pe.Time.Format("15:04:05"),   // format time as "HH:MM:SS"
		Alias: (*Alias)(&pe),
	})
}

// Administration requests

type SusReq struct {
	Target   string `json:"target"`
	Duration int64  `json:"duration"`
}

type LockReq struct {
	Title        string `json:"title"`
	MinimumGroup int64  `json:"minGroup"`
}

type RightsReq struct {
	Add    []string
	Remove []string
}
