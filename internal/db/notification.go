package db

import (
	"errors"
	"log"
)

func AddNotification(sort, date string, sender, receiver, post, parentPost int) error {
	db := GetDB()
	defer db.Close()

	if sender == receiver {
		return nil
	}

	if isSpam(sender, receiver, post) && sort == "report" {
		return errors.New("you have already reported this message")
	} else if isSpam(sender, receiver, post) {
		delSpamNotification(sender, receiver, post)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO notifications (sort, sender, receiver, post, parentPost, date) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sort, sender, receiver, post, parentPost, date)

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

func delSpamNotification(sender, receiver, post int) error {
	db := GetDB()
	defer db.Close()

	_, err := db.Exec(`DELETE FROM notifications WHERE sort IN ('LIKE', 'DISLIKE') AND sender = ? AND receiver = ? AND post = ?`, sender, receiver, post)
	if err != nil {
		log.Printf("Error when deleting session: %v", err)
		return err
	}
	return nil
}

func isSpam(sender, receiver, post int) bool {
	db := GetDB()
	defer db.Close()

	var isSpam bool
	err := db.QueryRow(`SELECT CASE WHEN EXISTS (SELECT * FROM notifications WHERE sort IN ('LIKE', 'DISLIKE', 'report') AND sender = ? AND receiver = ? AND post = ?) THEN CAST(1 AS BIT) ELSE CAST(0 AS BIT) END;`, sender, receiver, post).Scan(&isSpam)
	if err != nil {
		return isSpam
	}
	return isSpam
}
