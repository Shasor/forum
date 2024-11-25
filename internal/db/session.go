package db

import (
	"errors"
	"log"
)

func AddConnectedUser(userID int, sessionID string) error {
	// Ouvrir la connexion à la base de données
	db := GetDB()
	defer db.Close()

	// Démarrer une transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Préparer la requête d'insertion
	stmt, err := tx.Prepare("INSERT INTO sessions(connected_user, uuid) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Exécuter l'insertion
	_, err = stmt.Exec(userID, sessionID)
	if err != nil {
		tx.Rollback()
		return errors.New("la session existe déjà")
	}

	// Commit de la transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func IsUserConnected(userID int) (bool, error) {
	db := GetDB()
	defer db.Close()

	var exist bool
	err := db.QueryRow(`SELECT EXISTS (SELECT 1 FROM sessions WHERE connected_user = ? LIMIT 1)`, userID).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, nil
}

func DeleteConnectedUser(userID int) error {
	db := GetDB()
	defer db.Close()

	// Now delete the user
	_, err := db.Exec(`DELETE FROM sessions WHERE connected_user = ?`, userID)
	if err != nil {
		log.Printf("Error when deleting user: %v", err)
		return err
	}

	return nil
}
