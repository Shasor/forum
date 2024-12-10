package db

import (
	"errors"
	"fmt"
	"log"
)

func AddNotification(sort, date string, sender, receiver, post, parentPost int) error {
	db := GetDB()
	defer db.Close()

	if sender == receiver {
		return nil
	}

	if isSpam(sender, receiver, post) && sort == "request" {
		return errors.New("you've already asked to be a moderator")
	} else if isSpam(sender, receiver, post) && sort == "report" {
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
	err := db.QueryRow(`SELECT CASE WHEN EXISTS (SELECT * FROM notifications WHERE sort IN ('LIKE', 'DISLIKE', 'report', 'request') AND sender = ? AND receiver = ? AND post = ?) THEN CAST(1 AS BIT) ELSE CAST(0 AS BIT) END;`, sender, receiver, post).Scan(&isSpam)
	if err != nil {
		return isSpam
	}
	return isSpam
}

func FetchNotificationsByUserId(userID int) ([]Notification, error) {
	db := GetDB()
	defer db.Close()

	var query string
	if IsUserAdmin(userID) {
		query = `
		SELECT id, sort, sender, receiver, post, parentPost, readed, date 
		FROM notifications 
		WHERE (receiver = ? OR receiver = 0) AND readed = 0 ;
		`

	} else {
		query = `
		SELECT id, sort, sender, receiver, post, parentPost, readed, date 
		FROM notifications 
		WHERE receiver = ? AND readed = 0;
		`
	}
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Printf("Error executing the query: %v", err)
		return nil, err
	}

	defer rows.Close()
	var notifs []Notification
	for rows.Next() {
		var noti Notification
		err := rows.Scan(&noti.ID, &noti.Sort, &noti.Sender.ID, &noti.Receiver.ID, &noti.Post.ID, &noti.ParentID.ID, &noti.Readed, &noti.Date)
		if err != nil {
			fmt.Println("Error executing request : ", err)
		}
		noti.Sender, err = SelectUserByID(noti.Sender.ID)
		if err != nil && noti.Sender.ID != 0 {
			return nil, err
		}
		if noti.Post.ID != 0 {
			noti.Post, err = SelectPostByID(noti.Post.ID)
			if err != nil {
				continue
			}
		}

		if err = rows.Err(); err != nil {
			log.Printf("Error during row iteration: %v", err)
			return nil, err
		}
		notifs = append(notifs, noti)
	}
	return notifs, nil
}

func MarkAllNotificationsAsRead(userID int) error {
	db := GetDB()
	defer db.Close()

	query := "UPDATE notifications SET readed = 1 WHERE receiver = ? AND readed = 0"
	if IsUserAdmin(userID) {
		query = "UPDATE notifications SET readed = 1 WHERE (receiver = ? OR receiver = 0) AND readed = 0  "
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func CheckReport(id int) int {
	db := GetDB()
	defer db.Close()

	var senderID int
	err := db.QueryRow("SELECT sender FROM notifications WHERE post = ? AND sort = 'report'", id).Scan(&senderID)
	if err != nil {
		return 0
	}

	return senderID
}

func ReadNotification(sort string, sender, receiver, post int) error {
	db := GetDB()
	defer db.Close()

	_, err := db.Exec(`UPDATE notifications SET readed = 1 WHERE sort = ? AND sender = ? AND receiver = ? AND post = ?`, sort, sender, receiver, post)
	if err != nil {
		log.Printf("Error when set notification to readed: %v", err)
		return err
	}
	return nil
}
