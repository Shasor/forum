package db

import (
	"errors"
	"log"
)

func AddConnectedUser(userID int, sessionUUID string) error {
	db := GetDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO sessions(connected_user, uuid) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, sessionUUID)
	if err != nil {
		tx.Rollback()
		return errors.New("session already exists")
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func IsUserConnected(userID int) (bool, error) {
	db := GetDB()
	defer db.Close()

	var exists bool
	err := db.QueryRow(`SELECT EXISTS (SELECT 1 FROM sessions WHERE connected_user = ? LIMIT 1)`, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func DeleteConnectedUser(sessionUUID string) error {
	db := GetDB()
	defer db.Close()

	_, err := db.Exec(`DELETE FROM sessions WHERE uuid = ?`, sessionUUID)
	if err != nil {
		log.Printf("Error when deleting session: %v", err)
		return err
	}

	return nil
}

func GetUserIDBySessionUUID(sessionUUID string) (int, error) {
	db := GetDB()
	defer db.Close()

	var userID int
	err := db.QueryRow(`SELECT connected_user FROM sessions WHERE uuid = ?`, sessionUUID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func GetUUIDByUserID(id int) (string, error) {
	db := GetDB()
	defer db.Close()

	var uuid string
	err := db.QueryRow(`SELECT uuid FROM sessions WHERE connected_user = ?`, id).Scan(&uuid)
	if err != nil {
		return "", err
	}

	return uuid, nil
}
