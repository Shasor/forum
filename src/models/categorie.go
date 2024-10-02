package models

import (
	"database/sql"
	"log"
)

func CreateCategorie(db *sql.DB, Name string) error {
	// Insertion dans la base de données
	statement, err := db.Prepare("INSERT INTO Categorie (Name) VALUES (?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(Name)
	return err
}

func DeleteCategorie(db *sql.DB, CategorieID int) error {
	// Requête SQL pour supprimer l'utilisateur
	stmt, err := db.Prepare("DELETE FROM Post WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(CategorieID)
	if err != nil {
		log.Println("Erreur lors de la suppression de l'utilisateur :", err)
		return err
	}
	return nil
}
