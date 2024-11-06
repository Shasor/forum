package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func CreateCategory(name string) error {
	db := GetDB()
	defer db.Close()

	if CategoryExist(name) {
		return errors.New("the category already exists")
	}

	// Insertion dans la base de données
	statement, err := db.Prepare("INSERT INTO categories (name) VALUES (?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(name)
	if err != nil {
		return err
	}
	return nil
}

func SelectCategoryByName(name string) (Category, error) {
	db := GetDB()
	defer db.Close()

	var category Category
	err := db.QueryRow(`
	SELECT id, name
	FROM categories
	WHERE name = ?`,
		name).Scan(&category.ID, &category.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return Category{}, errors.New("user not found")
		}
		return Category{}, err
	}

	return category, nil
}

func CategoryExist(name string) bool {
	db := GetDB()
	defer db.Close()

	var exist bool
	err := db.QueryRow("SELECT EXISTS( SELECT 1 FROM categories WHERE name = ?)", name).Scan(&exist)
	if err != nil {
		return false
	}
	return exist
}

func GetCategoryNameByID(id int) string {
	db := GetDB()
	defer db.Close()

	var name string
	err := db.QueryRow("SELECT name FROM categories WHERE id = ?", id).Scan(&name)
	if err != nil {
		return ""
	}
	return name
}

func FetchCategories() []Category {
	db := GetDB()
	defer db.Close()

	query := `
        SELECT c.id, c.name
        FROM categories c;`

	rows, err := db.Query(query)
	if err != nil {
		// Gérer l'erreur, par exemple en la journalisant
		log.Printf("Erreur lors de l'exécution de la requête : %v", err)
		return nil
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name)
		fmt.Println(category)
		if err != nil {
			// Gérer l'erreur, par exemple en la journalisant
			log.Printf("Erreur lors du scan de la ligne : %v", err)
			continue
		}
		categories = append(categories, category)
	}
	if err = rows.Err(); err != nil {
		// Gérer l'erreur finale, si elle existe
		log.Printf("Erreur lors de l'itération sur les lignes : %v", err)
	}
	return categories
}
