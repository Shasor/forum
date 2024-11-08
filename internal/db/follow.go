package db

// AlreadyFollowingCategory vérifie si un utilisateur suit déjà une catégorie
func AlreadyFollowingCategory(categorieID, userID int) (bool, error) {
	db := GetDB()
	defer db.Close()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM follows WHERE user = ? AND category = ? LIMIT 1)`

	err := db.QueryRow(query, userID, categorieID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// StartFollowingCategory ajoute un suivi d'une catégorie par un utilisateur
func StartFollowingCategory(categorieID, userID int) error {
	db := GetDB()
	defer db.Close()

	// Préparation de la requête d'insertion
	statement, err := db.Prepare("INSERT INTO follows (category, user) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(categorieID, userID)
	return err
}

// StopFollowingCategory supprime le suivi d'une catégorie par un utilisateur
func StopFollowingCategory(categorieID, userID int) error {
	db := GetDB()
	defer db.Close()

	// Préparation de la requête de suppression
	statement, err := db.Prepare("DELETE FROM follows WHERE category = ? AND user = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(categorieID, userID)
	return err
}
