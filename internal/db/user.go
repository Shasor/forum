package db

import (
	"errors"
)

func CreateUser(role, username, email, picture, password string) error {
	// Ouvrir la connexion à la base de données
	db := GetDB()
	defer db.Close()

	// Démarrer une transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Préparer la requête d'insertion
	stmt, err := tx.Prepare("INSERT INTO users(role, username, email, picture, password) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Exécuter l'insertion
	_, err = stmt.Exec(role, username, email, picture, password)
	if err != nil {
		tx.Rollback()
		return errors.New("l'email ou le pseudo existe déjà")
	}

	// Commit de la transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
