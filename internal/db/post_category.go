package db

import (
	"log"
)

func LinkPostToCategory(postID int, categoryID int) error {
	db := GetDB()
	defer db.Close()

	// Démarrer une transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Préparer la requête d'insertion
	stmt, err := tx.Prepare("INSERT INTO post_category (post_id, category_id) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Exécuter l'insertion
	_, err = stmt.Exec(&postID, categoryID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit de la transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func GetPostCategories(postID int) ([]Category, error) {

	db := GetDB()
	defer db.Close()

	query := `
	SELECT category_id
	FROM post_category
	WHERE post_id = ?
	`
	var categories []Category

	rows, err := db.Query(query, postID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return categories, nil
	}

	defer rows.Close()

	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID); err != nil {
			log.Printf("Error scanning row: %v", err)
			return categories, nil
		}

		category.Name = GetCategoryNameByID(category.ID)

		categories = append(categories, category)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
	}

	//fmt.Println("Id du post : ", postID)
	//fmt.Println("Categories : ", categories)

	return categories, nil

}
