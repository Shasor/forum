package db

import (
	"database/sql"
	"log"
)

func UpdatePostReaction(userID int, postID int, reaction string) error {
	db := GetDB()
	defer db.Close()

	var existingReaction string
	err := db.QueryRow(`SELECT value FROM reactions WHERE sender = ? AND post = ?`, userID, postID).Scan(&existingReaction)

	if err == sql.ErrNoRows {
		// No previous reaction: Insert new
		_, err = db.Exec(`INSERT INTO reactions (sender, post, value) VALUES (?, ?, ?)`, userID, postID, reaction)
		if err != nil {
			log.Println("Error inserting new reaction:", err)
			return err
		}
	} else if err != nil {
		log.Println("Error checking existing reaction:", err)
		return err
	} else {
		// If user reacted the same way before, delete it (toggle reaction off)
		if existingReaction == reaction {
			_, err = db.Exec(`DELETE FROM reactions WHERE sender = ? AND post = ?`, userID, postID)
			if err != nil {
				log.Println("Error removing existing reaction:", err)
				return err
			}
		} else {
			// Otherwise, update the reaction
			_, err = db.Exec(`UPDATE reactions SET value = ? WHERE sender = ? AND post = ?`, reaction, userID, postID)
			if err != nil {
				log.Println("Error updating existing reaction:", err)
				return err
			}
		}
	}
	return nil
}

func GetPostReactions(postID int) (int, int, error) {
	db := GetDB()
	defer db.Close()

	var likes, dislikes int
	err := db.QueryRow(`SELECT COUNT(*) FROM reactions WHERE post = ? AND value = 'LIKE'`, postID).Scan(&likes)
	if err != nil {
		return 0, 0, err
	}

	err = db.QueryRow(`SELECT COUNT(*) FROM reactions WHERE post = ? AND value = 'DISLIKE'`, postID).Scan(&dislikes)
	if err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}
