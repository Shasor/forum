package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func UpdatePostReaction(userID int, postID int, reaction string) error {
	db := GetDB()
	defer db.Close()

	date := fmt.Sprintf("%02d:%02d | %02d/%02d/%d", time.Now().Hour(), time.Now().Minute(), time.Now().Day(), time.Now().Month(), time.Now().Year())
	post, _ := SelectPostByID(postID)

	var existingReaction string
	err := db.QueryRow(`SELECT value FROM reactions WHERE user = ? AND post = ?`, userID, postID).Scan(&existingReaction)

	if err == sql.ErrNoRows {
		// No previous reaction: Insert new
		_, err = db.Exec(`INSERT INTO reactions (user, post, value) VALUES (?, ?, ?)`, userID, postID, reaction)
		addActivity(userID, postID, reaction)
		addNotification(reaction, date, userID, post.Sender.ID, post.ID, 0)
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
			_, err = db.Exec(`DELETE FROM reactions WHERE user = ? AND post = ?`, userID, postID)
			delActivity(userID, postID, reaction)
			if err != nil {
				log.Println("Error removing existing reaction:", err)
				return err
			}
		} else {
			// Otherwise, update the reaction
			_, err = db.Exec(`UPDATE reactions SET value = ? WHERE user = ? AND post = ?`, reaction, userID, postID)
			if reaction == "LIKE" {
				delActivity(userID, postID, "DISLIKE")
			} else {
				delActivity(userID, postID, "LIKE")
			}
			addActivity(userID, postID, reaction)
			addNotification(reaction, date, userID, post.Sender.ID, post.ID, 0)
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

func GetReactionsByPostID(id int) []Reaction {
	db := GetDB()
	defer db.Close()

	query := `
        SELECT r.id, r.user,  r.value
        FROM reactions r 
        WHERE r.post = ?;`

	// Execute the query
	rows, err := db.Query(query, id)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil
	}
	defer rows.Close()

	// Slice to hold the posts
	var reactions []Reaction

	for rows.Next() {
		var reaction Reaction
		var senderID int
		// Scan each row's values, including placeholders for missing user info
		if err := rows.Scan(&reaction.ID, &senderID, &reaction.Value); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil
		}

		reaction.Sender, _ = SelectUserById(senderID)

		reactions = append(reactions, reaction)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
		return nil
	}

	// Return the list of liked posts
	return reactions
}
