package models

import (
	"database/sql"
	"errors"
	"log"
	"login/src/database"
)

func CategorieExist(db *sql.DB, Name string) bool {

	var exists bool

	err := db.QueryRow("SELECT EXISTS( SELECT 1 FROM Categories WHERE Name = ?)", Name).Scan(&exists)
	if err != nil {
		return false
	}
	return exists

}

func CreateCategorie(db *sql.DB, Name string) error {
	// Insertion dans la base de données

	statement, err := db.Prepare("INSERT INTO Categories (Name) VALUES (?)")
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
	err := db.QueryRow("SELECT id, Name FROM Categories WHERE id = ?", CateID).Scan(
		&categorie.CategorieID, &categorie.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return categorie, errors.New("catégorie non trouvé")
		}
		return categorie, err
	}
	return categorie, nil
}

// GetCategorieIDByName retourne l'ID d'une catégorie à partir de son nom
func GetCategorieIDByName(db *sql.DB, name string) (int, error) {
	var id int
	// Exécution de la requête pour récupérer l'ID en fonction du nom
	err := db.QueryRow("SELECT id FROM Categories WHERE Name = ?", name).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("catégorie non trouvée")
		}
		return 0, err
	}
	return id, nil
}

// FetchCategories fetches a limited number of users from the database, excluding the current user's account
func FetchCategories(db *sql.DB) ([]database.Categorie, error) {
	// Use a parameterized query to prevent SQL injection
	query := "SELECT * FROM Categories"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []database.Categorie
	for rows.Next() {
		var categorie database.Categorie
		err := rows.Scan(&categorie.CategorieID, &categorie.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, categorie)
	}

	return categories, nil
}
