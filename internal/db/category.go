package db

import (
	"database/sql"
	"errors"
)

func CreateCategorie(name string) error {
	db := GetDB()
	defer db.Close()

	if CategorieExist(name) {
		return errors.New("the category already exists")
	}

	// Insertion dans la base de données
	statement, err := db.Prepare("INSERT INTO categories (name) VALUES (?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(name)
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

func CategorieExist(name string) bool {
	db := GetDB()
	defer db.Close()

	var exist bool
	err := db.QueryRow("SELECT EXISTS( SELECT 1 FROM categories WHERE name = ?)", name).Scan(&exist)
	if err != nil {
		return false
	}
	return exist
}
