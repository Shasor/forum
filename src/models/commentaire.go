package models

import (
	"database/sql"
	"errors"
	"log"
	"login/src/database"
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

// Récupère un commentaire par son id
func SelectCommentaire(db *sql.DB, CommentID int) (database.Commentaire, error) {
	var commentaire database.Commentaire
	err := db.QueryRow("SELECT id, SenderID, PostID, Like, Dislike, Date, Content FROM Commentaire WHERE id = ?", CommentID).Scan(
		&commentaire.CommentaireID, &commentaire.SenderID, &commentaire.PostID, &commentaire.Like, &commentaire.Dislike, &commentaire.Date, &commentaire.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return commentaire, errors.New("commentaire non trouvé")
		}
		return commentaire, err
	}
	return commentaire, nil
}
