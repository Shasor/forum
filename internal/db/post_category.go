package db

func LinkPostToCategory(postID int, category Category) error {

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
	_, err = stmt.Exec(&postID, category.ID)
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


    return nil
}
