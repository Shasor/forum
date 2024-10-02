package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreatePost(db *sql.DB, Title, Content, Date string, Sender int, Image, Like, Dislike string) error {
	// Insertion dans la base de données
	statement, err := db.Prepare("INSERT INTO Users (Title, Content, Date, Sender, Image, Like, Dislike) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(Title, Content, Date, Sender, Image, Like, Dislike)
	return err
}

func DeletePost(db *sql.DB, PostID int) error {
	// Requête SQL pour supprimer l'utilisateur
	stmt, err := db.Prepare("DELETE FROM Post WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(PostID)
	if err != nil {
		log.Println("Erreur lors de la suppression de l'utilisateur :", err)
		return err
	}
	return nil
}
