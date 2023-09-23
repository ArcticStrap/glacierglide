package data

import (
	"context"
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
func (db *PostgresBase) CreatePage(v *Page) error {
	// Execute create request
	_, err := db.conn.Exec(context.Background(), "INSERT INTO pages (title, content) VALUES ($1, $2)", v.Title, v.Content)
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

func (db *PostgresBase) UpdatePage(v *Page) error {
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
