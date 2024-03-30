package data

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/ArcticStrap/glacierglide/utils/iputils"
	"github.com/ArcticStrap/glacierglide/utils/pageutils/pnamespace"
	"github.com/ArcticStrap/glacierglide/utils/userutils"
	"github.com/ArcticStrap/glacierglide/wikiconfig"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// Created for use in future updates when we plan on implementing different databses
type Datastore interface {
	// Initalizers / Cleaners
	Close()
	CreateTables() error

	// Metadata
	EngineName() string
	Version() string

	// CRUD functions for page
	CreatePage(v *Page) error
	ReadPage(title string) (*Page, error)
	UpdatePage(p *Page, editor string) error
	DeletePage(title string) error

	// Small gets
	GetIdFromPageTitle(title string) (*int, error)
	GetUserIdFromName(username string) (*int64, error)
	GetUsernameFromId(id int64) (string, error)
	FetchPageHistory(title string) ([]PageEdit, error)

	// CRUD for user
	CreateUser(username string, password string) (*User, error)
	GetUser(username string) (*User, error)
	UpdateUser(u *User, newName string, newPass string) error
	DeleteUser(username string) error

	GetUserGroups(username string) ([]string, error)
	EditUserGroups(username string, changes RightsReq) error

	// CRUD for page edit
	CreatePageEdit(p *Page, editor string) error
	ReadPageEdit(id int64) (*PageEdit, error)
	DeletePageEdit(id int64) error

	// Search operations
	SearchPagesFromTitle(title string) ([]Page, error)
	SearchPagesFromTitlePrefix(prefix string, limit int) ([]Page, error)
	SearchPagesContainingTitle(title string, limit int) ([]Page, error)

	// Moderation actions
	GetLockStatus(p string) (int, error)
	LockPage(p string, min_group int) error
	UnlockPage(p string, min_group int) error

	IsSuspended(u string) (bool, error)
	SuspendUser(u string, duration int64) error
}

// Datastore implementation for the Postgres database.

type PostgresBase struct {
	conn *pgx.Conn
}

// Postgres connections
func ConnectToPostgresDatabase() (PostgresBase, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("PAGEDATAURL"))
	if err != nil {
		return PostgresBase{}, err
	}
	return PostgresBase{conn: conn}, nil
}

func (db *PostgresBase) CreateTables() error {
	// Create the tables for pages, users, page history, and diffs if they don't exist
	_, err := db.conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS pages (
			page_id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
      namespace INT NOT NULL,
      mp_type INT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS page_edits (
			edit_id SERIAL PRIMARY KEY,
			page_id INT REFERENCES pages(page_id),
			change_date DATE NOT NULL,
			change_time TIME NOT NULL,
			editor TEXT,
			description TEXT,
			content TEXT
		);
		CREATE TABLE IF NOT EXISTS users (
			user_id SERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			creation_date TIMESTAMP
		);
    CREATE TABLE IF NOT EXISTS user_groups (
      user_id INT NOT NULL REFERENCES users(user_id),
      u_group INT NOT NULL,
      PRIMARY KEY (user_id,u_group)
    );
    CREATE TABLE IF NOT EXISTS suspensions (
      sus_id SERIAL PRIMARY KEY,
      sus_target TEXT NOT NULL,
      sus_ends INT NOT NULL
    );
    CREATE TABLE IF NOT EXISTS locks (
      lock_id SERIAL PRIMARY KEY,
      page_id INT REFERENCES pages(page_id),
      min_group INT
    );
	`)
	if err != nil {
		return err
	}
	return nil
}

func (db *PostgresBase) Close() {
	db.conn.Close(context.Background())
}

// Metadata
func (db *PostgresBase) EngineName() string {
	return "PostgreSQL"
}

func (db *PostgresBase) Version() string {
	var serverVersion string
	err := db.conn.QueryRow(context.Background(), "SHOW server_version").Scan(&serverVersion)
	if err != nil {
		return ""
	}
	return serverVersion
}

// Fetches page id from page title
func (db *PostgresBase) GetIdFromPageTitle(title string) (*int, error) {
	// Check for namespace
	ns := 0
	if sTitle := strings.Split(title, ":"); len(sTitle) == 2 {
		ns = pnamespace.NumberFromNamespace(sTitle[0])
		title = sTitle[1]
	}

	var id int
	err := db.conn.QueryRow(context.Background(), "SELECT page_id FROM pages WHERE title=$1 AND namespace=$2", title, ns).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// CRUD functions for Page datatype

// Logs a page change into a diff.
func (db *PostgresBase) CreatePageEdit(p *Page, editor string) error {
	// SQL query to insert a new page_edits row
	query := "INSERT INTO page_edits (page_id,change_date,change_time,editor,description,content) VALUES ($1, $2, $3, $4, $5, $6)"

	// Get page id
	pId, err := db.GetIdFromPageTitle(p.Title)
	if err != nil {
		return err
	}

	// Execute create request
	_, err = db.conn.Exec(context.Background(), query, pId, time.Now().UTC().Format("2006-01-02"), time.Now().UTC().Format("15:04:05"), editor, "edit page", p.Content)
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresBase) ReadPageEdit(id int64) (*PageEdit, error) {
	// SQL query to fetch the page by ID
	query := `SELECT page_id,change_date,change_time,editor,description,content FROM page_edits WHERE edit_id=$1`

	var pageEdit PageEdit

	// Execute the query and scan the result into the Page struct
	if err := db.conn.QueryRow(context.Background(), query, id).Scan(&pageEdit.PageId, &pageEdit.Date, &pageEdit.Time, &pageEdit.UserId, &pageEdit.Description, &pageEdit.Content); err != nil {
		return nil, err
	}

	pageEdit.EditId = id

	return &pageEdit, nil
}

func (db *PostgresBase) DeletePageEdit(id int64) error {
	// Execute delete request for page history
	_, err := db.conn.Exec(context.Background(), "DELETE FROM page_edits WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresBase) ReadPage(title string) (*Page, error) {
	// SQL query to fetch the page by ID
	query := `SELECT title, content, namespace, mp_type FROM pages WHERE page_id=$1`

	pageID, err := db.GetIdFromPageTitle(title)
	if err != nil {
		return nil, err
	}

	var page Page

	// Execute the query and scan the result into the Page struct
	if err := db.conn.QueryRow(context.Background(), query, pageID).Scan(&page.Title, &page.Content, &page.Namespace, &page.MPType); err != nil {
		return nil, err
	}

	return &page, nil
}

func (db *PostgresBase) CreatePage(v *Page) error {
	// Execute create request
	_, err := db.conn.Exec(context.Background(), "INSERT INTO pages (title, content, namespace, mp_type) VALUES ($1, $2, $3, $4)", strings.ToLower(v.Title), v.Content, v.Namespace, v.MPType)
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresBase) UpdatePage(p *Page, editor string) error {

	// Fetch page id from title
	id, err := db.GetIdFromPageTitle(p.Title)
	if err != nil {
		return err
	}

	// Execute create request
	_, err = db.conn.Exec(context.Background(), "UPDATE pages SET content = ($1) WHERE page_id = ($2)", p.Content, id)
	if err != nil {
		return err
	}

	// Log update
	if err = db.CreatePageEdit(p, editor); err != nil {
		return err
	}

	return nil
}

func (db *PostgresBase) DeletePage(title string) error {
	// Fetch page id from title
	id, err := db.GetIdFromPageTitle(title)
	if err != nil {
		return err
	}

	// Execute delete request
	_, err = db.conn.Exec(context.Background(), "DELETE FROM pages WHERE page_id=$1", *id)
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresBase) FetchPageHistory(title string) ([]PageEdit, error) {
	// History var
	pageHistory := []PageEdit{}

	// Get page info
	pgID, err := db.GetIdFromPageTitle(title)
	if err != nil {
		return nil, err
	}

	// Query rows
	rows, err := db.conn.Query(context.Background(), "SELECT * FROM page_edits WHERE page_id=$1", *pgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ok := true
	for rows.Next() {
		var pd PageEdit
		err = rows.Scan(&pd.EditId, &pd.PageId, &pd.Date, &pd.Time, &pd.UserId, &pd.Description, &pd.Content)
		if err != nil {
			ok = false
			break
		}

		pageHistory = append(pageHistory, pd)
	}
	if !ok {
		return nil, err
	}

	return pageHistory, nil
}

// CRUD functions for user accounts

func (db *PostgresBase) GetUserIdFromName(username string) (*int64, error) {
	var id int64
	err := db.conn.QueryRow(context.Background(), "SELECT user_id FROM users WHERE username=$1", username).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (db *PostgresBase) GetUsernameFromId(id int64) (string, error) {
	var username string
	err := db.conn.QueryRow(context.Background(), "SELECT username FROM users WHERE user_id=$1", id).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func (db *PostgresBase) GetUser(username string) (*User, error) {
	// SQL query to fetch the user by ID
	query := `SELECT username, password, creation_date FROM users WHERE user_id=$1`

	userID, err := db.GetUserIdFromName(username)
	if err != nil {
		return nil, err
	}

	var u User

	// Execute the query and scan the result into the Page struct
	if err := db.conn.QueryRow(context.Background(), query, userID).Scan(&u.Username, &u.Password, &u.CreationDate); err != nil {
		return nil, err
	}

	u.UserId = int64(*userID)

	return &u, nil
}

func (db *PostgresBase) CreateUser(username, password string) (*User, error) {
	// Set up struct
	newUser := &User{
		Username:     username,
		Password:     password,
		CreationDate: time.Now().UTC(),
	}

	// Hash the password
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser.Password = string(hashPass)

	// Execute create request
	_, err = db.conn.Exec(context.Background(), "INSERT INTO users (username, password, creation_date) VALUES ($1, $2, $3)", username, newUser.Password, newUser.CreationDate)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (db *PostgresBase) UpdateUser(u *User, newName, newPass string) error {
	// Fetch page id from name
	id, err := db.GetUserIdFromName(u.Username)
	if err != nil {
		return err
	}

	// Execute create request
	_, err = db.conn.Exec(context.Background(), "UPDATE users SET username = ($1), SET password = ($2) WHERE user_id = ($3)", newName, newPass, id)
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresBase) DeleteUser(username string) error {
	// Fetch user id from name
	id, err := db.GetUserIdFromName(username)
	if err != nil {
		return err
	}

	// Execute delete request
	_, err = db.conn.Exec(context.Background(), "DELETE FROM users WHERE user_id=$1", *id)
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresBase) GetUserGroups(username string) ([]string, error) {
	if iputils.NameIsIP(username) {
		return []string{"*"}, nil
	}

	// Fetch user id from name
	id, err := db.GetUserIdFromName(username)
	if err != nil {
		return nil, err
	}

	var groups = []string{"*", wikiconfig.DefaultLoginGroup}

	// Execute request
	rows, err := db.conn.Query(context.Background(), "SELECT u_group FROM user_groups WHERE user_id=$1", *id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ok := true
	for rows.Next() {
		var gID int
		err = rows.Scan(&gID)
		if err != nil {
			ok = false
			break
		}
		groups = append(groups, userutils.GroupFromIndex(gID))
	}
	if !ok {
		return nil, err
	}

	return groups, nil
}

func (db *PostgresBase) EditUserGroups(username string, changes RightsReq) error {
	tx, err := db.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Add users to groups
	for _, group := range changes.Add {
		_, err := tx.Exec(context.Background(), "INSERT INTO user_groups (user_id, u_group) VALUES ((SELECT user_id FROM users WHERE username = $1), $2)", username, group)
		if err != nil {
			return err
		}
	}

	// Remove users from groups
	for _, group := range changes.Remove {
		_, err := tx.Exec(context.Background(), "DELETE FROM user_groups WHERE user_id = (SELECT user_id FROM users WHERE username = $1) AND u_group = $2", username, group)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// Search operatons
func (db *PostgresBase) SearchPagesFromTitle(title string) ([]Page, error) {
	var results []Page

	// Query
	rows, err := db.conn.Query(context.Background(), "SELECT * FROM pages WHERE title=$1", title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ok := true
	for rows.Next() {
		var rPage Page
		err = rows.Scan(nil, &rPage.Title, &rPage.Content, &rPage.Namespace, &rPage.MPType)
		if err != nil {
			ok = false
			break
		}
		results = append(results, rPage)
	}
	if !ok {
		return nil, err
	}

	return results, nil
}

func (db *PostgresBase) SearchPagesFromTitlePrefix(prefix string, limit int) ([]Page, error) {
	var results []Page

	// Query
	rows, err := db.conn.Query(context.Background(), "SELECT * FROM pages WHERE title LIKE $1 || '%' LIMIT $2", prefix, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ok := true
	for rows.Next() {
		var rPage Page
		err = rows.Scan(nil, &rPage.Title, &rPage.Content, &rPage.Namespace, &rPage.MPType)
		if err != nil {
			ok = false
			break
		}
		results = append(results, rPage)
	}
	if !ok {
		println(err.Error())
		return nil, err
	}

	return results, nil
}

func (db *PostgresBase) SearchPagesContainingTitle(title string, limit int) ([]Page, error) {
	results := []Page{}

	// Query
	rows, err := db.conn.Query(context.Background(), "SELECT * FROM pages WHERE title LIKE '%' || $1 || '%' LIMIT $2", title, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ok := true
	for rows.Next() {
		var rPage Page
		err = rows.Scan(nil, &rPage.Title, &rPage.Content, &rPage.Namespace, &rPage.MPType)
		if err != nil {
			ok = false
			break
		}
		results = append(results, rPage)
	}
	if !ok {
		return nil, err
	}

	return results, nil
}

// Moderation actions

func (db *PostgresBase) GetLockStatus(p string) (int, error) {
	// Fetch page id
	id, err := db.GetIdFromPageTitle(p)
	if err != nil {
		return 0, err
	}

	// Query rows
	rows, err := db.conn.Query(context.Background(), "SELECT min_group FROM locks WHERE page_id=$1", *id)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	min_group := 0

	ok := true
	for rows.Next() {
		err = rows.Scan(&min_group)
		if err != nil {
			ok = false
			break
		}
	}
	if !ok {
		return 0, err
	}

	return min_group, nil
}

func (db *PostgresBase) LockPage(p string, min_group int) error {
	// Fetch page id
	id, err := db.GetIdFromPageTitle(p)
	if err != nil {
		return err
	}

	// Execute request
	_, err = db.conn.Exec(context.Background(), "INSERT INTO locks (page_id,min_group) VALUES ($1,$2)", *id, min_group)
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresBase) UnlockPage(p string, min_group int) error {
	// Fetch page id
	id, err := db.GetIdFromPageTitle(p)
	if err != nil {
		return err
	}

	// Execute request
	_, err = db.conn.Exec(context.Background(), "DELETE FROM locks WHERE page_id=$1", *id)
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresBase) IsSuspended(u string) (bool, error) {
	rows, err := db.conn.Query(context.Background(), "SELECT sus_ends FROM suspensions WHERE sus_target=$1", u)
	if err != nil {
		return false, err
	}

	var sDuration int64

	ok := true
	for rows.Next() {
		err = rows.Scan(&sDuration)
		if err != nil {
			ok = false
			break
		}
	}
	if !ok {
		return false, err
	}

	var blocked = time.Now().Unix() < sDuration
	// Remove block from database if done
	if !blocked {
		_, err := db.conn.Exec(context.Background(), "DELETE FROM suspensions WHERE sus_target=$1", u)
		if err != nil {
			return false, err
		}
	}

	return blocked, nil
}

func (db *PostgresBase) SuspendUser(u string, duration int64) error {
	// Execute request
	_, err := db.conn.Exec(context.Background(), "INSERT INTO suspensions (sus_target,sus_ends) VALUES ($1,$2)", u, time.Now().Unix()+duration)
	if err != nil {
		return err
	}

	return nil
}
