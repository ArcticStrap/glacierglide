package data

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectToDataBase() (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), os.Getenv("PAGEDATAURL"))
}

func CreateTables(conn *pgx.Conn) error {
	// Create the tables for pages, users, page history, and diffs if they don't exist
	_, err := conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS pages (
			page_id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT not null
		);
		CREATE TABLE IF NOT EXISTS users (
			user_id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			creation_date TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS page_history (
			id SERIAL PRIMARY KEY,
			page_id INT REFERENCES pages(page_id),
			change_date DATE NOT NULL,
			change_time TIME NOT NULL,
			editor_id INT REFERENCES users(user_id),
			description TEXT,
			diff_id UUID UNIQUE
		);
		CREATE TABLE IF NOT EXISTS diffs (
			diff_id UUID PRIMARY KEY,
			diff_content TEXT NOT NULL -- this could be in a unified diff format or any format you choose
		);		
	`)
	if err != nil {
		return err
	}
	return nil
}

func CloseDataBase(conn *pgx.Conn) {
	conn.Close(context.Background())
}

func InsertPage(conn *pgx.Conn, v *Page) error {
	// Execute create request
	_, err := conn.Exec(context.Background(), "INSERT INTO pages (title, content) VALUES ($1, $2)", v.Title, v.Content)
	if err != nil {
		return err
	}

	return nil
}

func DeletePage(conn *pgx.Conn, title string) error {
	// Fetch page id from title
	var id int
	err := conn.QueryRow(context.Background(), "SELECT page_id FROM pages WHERE title=$1", title).Scan(&id)
	if err != nil {
		return err
	}

	// Execute delete request
	_, err = conn.Exec(context.Background(), "DELETE FROM pages WHERE page_id=$1", id)
	if err != nil {
		return err
	}

	return nil
}
