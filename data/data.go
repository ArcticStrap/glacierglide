package data

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

type PostgresBase struct {
	conn *pgx.Conn
}

// (*pgx.Conn, error)
func ConnectToDatabase() (PostgresBase, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("PAGEDATAURL"))
	if err != nil {
		return PostgresBase{}, err
	}
	return PostgresBase{conn: conn}, nil
}

func (b *PostgresBase) CreateTables() error {
	// Create the tables for pages, users, page history, and diffs if they don't exist
	_, err := b.conn.Exec(context.Background(), `
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
			anon BOOLEAN DEFAULT FALSE,
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

func (b *PostgresBase) Close() {
	b.conn.Close(context.Background())
}
