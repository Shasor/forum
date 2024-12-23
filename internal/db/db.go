package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func GetDB() *sql.DB {
	db, err := sql.Open("sqlite3", "internal/db/database.db")
	if err != nil {
		log.Fatal(err)
	}

	// SQL statement to create the users table
	tables := `CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			provider TEXT,
			role TEXT NOT NULL,
			username TEXT UNIQUE,
			email TEXT UNIQUE,
			picture TEXT,
			password TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		); CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);  CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category INTEGER NOT NULL,
			sender INTEGER NOT NULL,
			parent_id INTEGER DEFAULT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			picture LONGTEXT,
			date TEXT NOT NULL,
			FOREIGN KEY (category) REFERENCES categories(id),
			FOREIGN KEY (sender) REFERENCES users(id)
			FOREIGN KEY (parent_id) REFERENCES posts(id)
		); CREATE TABLE IF NOT EXISTS post_category (
            post_id INTEGER,
            category_id INTEGER,
            FOREIGN KEY (post_id) REFERENCES post(id),
            FOREIGN KEY (category_id) REFERENCES categories(id)
		); CREATE TABLE IF NOT EXISTS reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user INTEGER NOT NULL,
			post INTEGER NOT NULL,
			value TEXT NOT NULL,
			FOREIGN KEY (user) REFERENCES users(id),
			FOREIGN KEY (post) REFERENCES posts(id)
		); CREATE TABLE IF NOT EXISTS follows (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		category INTEGER NOT NULL,
    		user INTEGER NOT NULL,
    		FOREIGN KEY (category) REFERENCES categories(id),
    		FOREIGN KEY (user) REFERENCES users(id)
		); CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			connected_user INTEGER NOT NULL,
			uuid STRING NOT NULL,
			FOREIGN KEY (connected_user) REFERENCES users(id)
		); CREATE TABLE IF NOT EXISTS activity (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		user INTEGER NOT NULL,
    		post INTEGER NOT NULL,
    		action TEXT NOT NULL,
    		FOREIGN KEY (user) REFERENCES users(id),
    		FOREIGN KEY (post) REFERENCES posts(id)
		); CREATE TABLE IF NOT EXISTS notifications (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		sort TEXT NOT NULL,
    		sender INTEGER NOT NULL,
    		receiver INTEGER NOT NULL,
    		post INTEGER,
			parentPost INTEGER,
			readed INTERGER NOT NULL DEFAULT 0,
			date TEXT NOT NULL,
    		FOREIGN KEY (sender) REFERENCES users(id),
    		FOREIGN KEY (receiver) REFERENCES users(id),
    		FOREIGN KEY (post) REFERENCES posts(id)
    		FOREIGN KEY (parentPost) REFERENCES posts(id)
		)`

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Execute the SQL statements within the transaction
	_, err = tx.Exec(tables)
	if err != nil {
		// If there's an error, roll back the transaction
		tx.Rollback()
		log.Fatal(err)
	}

	// If everything is okay, commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
