package models

import (
	"database/sql"
	"errors"
	"log"
	"login/src/database"
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

// Récupère une catégorie par son id
func SelectCategorie(db *sql.DB, CateID int) (database.Categorie, error) {
	var categorie database.Categorie
	err := db.QueryRow("SELECT id, Name FROM Categorie WHERE id = ?", CateID).Scan(
		&categorie.CategorieID, &categorie.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return categorie, errors.New("catégorie non trouvé")
		}
		return categorie, err
	}
	return categorie, nil
}
