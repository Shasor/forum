package models

import (
	"database/sql"
	"log"
)

func CreateCommentaire(db *sql.DB, SenderID, PostID int, Like, Dislike, Date, Content string) error {
	// Insertion dans la base de données
	statement, err := db.Prepare("INSERT INTO Commentaire (SenderID, PostID, Like, Dislike, Date, Content) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(SenderID, PostID, Like, Dislike, Date, Content)
	return err
}

func DeleteCommentaire(db *sql.DB, CommentaireID int) error {
	// Requête SQL pour supprimer l'utilisateur
	stmt, err := db.Prepare("DELETE FROM Post WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(CommentaireID)
	if err != nil {
		log.Println("Erreur lors de la suppression de l'utilisateur :", err)
		return err
	}
	return nil
}
