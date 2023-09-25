package data

import (
	"context"
	"time"
)

// Created for use in future updates when we plan on implementing different databses (e.g. MariaDB/MySQL and MongoDB)
/* type Datastore interface {
	Insert(any)
	Fetch(any)
	Delete(any)
	Update(any)
} */

func (db *PostgresBase) GetIdFromPageTitle(title string) (*int, error) {
	var id int
	err := db.conn.QueryRow(context.Background(), "SELECT page_id FROM pages WHERE title=$1", title).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// CRUD functions for Page datatype

func (db *PostgresBase) CreatePageDiff(p *Page) error {
	// SQL query to insert a new page_diffs row
	query := "INSERT INTO page_diffs (page_id,change_date,change_time,editor_id,anon,description,content) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	// Get page id
	pId, err := db.GetIdFromPageTitle(p.Title)
	if err != nil {
		return err
	}

	// Execute create request
	_, err = db.conn.Exec(context.Background(), query, pId, time.Now().UTC().Format("2006-01-02"), time.Now().UTC().Format("15:04:05"), nil, false, "edit page", p.Content)
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresBase) ReadPageDiff(id int64) (*PageDiff, error) {
	// SQL query to fetch the page by ID
	query := `SELECT page_id,change_date,change_time,editor_id,anon,description,diff_id FROM page_diffs WHERE id=$1`

	var pageDiff PageDiff

	// Execute the query and scan the result into the Page struct
	if err := db.conn.QueryRow(context.Background(), query, id).Scan(&pageDiff.PageId, &pageDiff.Date, &pageDiff.Time, &pageDiff.UserId, &pageDiff.Anon, &pageDiff.Description, &pageDiff.Content); err != nil {
		return nil, err
	}

	return &pageDiff, nil
}

func (db *PostgresBase) DeletePageDiff(id int64) error {
	// Execute delete request for page history
	_, err := db.conn.Exec(context.Background(), "DELETE FROM page_diffs WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresBase) ReadPage(title string) (*Page, error) {
	// SQL query to fetch the page by ID
	query := `SELECT id, title, content FROM pages WHERE id=$1`

	pageID, err := db.GetIdFromPageTitle(title)
	if err != nil {
		return nil, err
	}

	var page Page

	// Execute the query and scan the result into the Page struct
	if err := db.conn.QueryRow(context.Background(), query, pageID).Scan(&page.Title, &page.Content); err != nil {
		return nil, err
	}

	return &page, nil
}

func (db *PostgresBase) CreatePage(v *Page) error {
	// Execute create request
	_, err := db.conn.Exec(context.Background(), "INSERT INTO pages (title, content) VALUES ($1, $2)", v.Title, v.Content)
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresBase) UpdatePage(p *Page) error {

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
	if err = db.CreatePageDiff(p); err != nil {
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

// CRUD functions for user accounts

func (db *PostgresBase) GetUserIdFromName(username string) (*int, error) {
	var id int
	err := db.conn.QueryRow(context.Background(), "SELECT user_id FROM users WHERE username=$1", username).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (db *PostgresBase) GetUser(username string) (*User, error) {
	// SQL query to fetch the user by ID
	query := `SELECT username, password, creation_date FROM users WHERE user_id=$1`

	pageID, err := db.GetUserIdFromName(username)
	if err != nil {
		return nil, err
	}

	var u User

	// Execute the query and scan the result into the Page struct
	if err := db.conn.QueryRow(context.Background(), query, pageID).Scan(&u.Username, u.Password, u.CreationDate); err != nil {
		return nil, err
	}

	return &u, nil
}

func (db *PostgresBase) CreateUser(username, password string) (*User, error) {
	// Set up struct
	newUser := &User{
		Username:     username,
		Password:     password,
		CreationDate: time.Now().UTC(),
	}

	// Execute create request
	_, err := db.conn.Exec(context.Background(), "INSERT INTO users (username, password, creation_date) VALUES ($1, $2, $3)", username, password, newUser.CreationDate)
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
