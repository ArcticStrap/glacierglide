package data

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

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
			password TEXT NOT NULL,
			creation_date TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS page_diffs (
			diff_id SERIAL PRIMARY KEY,
			page_id INT REFERENCES pages(page_id),
			change_date DATE NOT NULL,
			change_time TIME NOT NULL,
			editor_id INT REFERENCES users(user_id),
			anon BOOLEAN DEFAULT FALSE,
			description TEXT,
			content TEXT
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
