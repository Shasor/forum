package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	filepath = "db/database.db"
)

// Initialise la base de données
func GetDB() (*sql.DB, error) {
	err := os.MkdirAll("db", 0755)
	if err != nil {
		return nil, err
	}
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
