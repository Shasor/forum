package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Initialise la base de données
func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	// Tester la connexion
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// Crée la table des utilisateurs si elle n'existe pas encore
func CreateTable(db *sql.DB) {
	statement, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    )`)
	statement.Exec()
}
