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
