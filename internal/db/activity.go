package db

import "log"

func addActivity(userID, postID int, action string) error {
	db := GetDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO activity (user, post, action) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, postID, action)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func delActivity(userID, postID int, action string) error {
	db := GetDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("DELETE FROM activity WHERE user = ? AND post = ? AND action = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, postID, action)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func GetUserActivitiesByID(userID int) []Activity {
	db := GetDB()
	defer db.Close()

	query := `
			SELECT a.id, a.post, a.action
			FROM activity a
			WHERE a.user = ?
			ORDER BY a.id DESC;`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil
	}
	defer rows.Close()

	var activities []Activity
	for rows.Next() {
		var activity Activity
		var postID int
		// Scan each row's values, including placeholders for missing user info
		if err := rows.Scan(&activity.ID, &postID, &activity.Action); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil
		}

		activity.User, _ = SelectUserById(userID)
		activity.Post, _ = SelectPostByID(postID)

		activities = append(activities, activity)
	}

	return activities
}
