package db

import (
	"database/sql"
	"errors"
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
