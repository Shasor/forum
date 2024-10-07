package models

import (
	"database/sql"
	"errors"
	"log"
	"login/src/database"

	_ "github.com/mattn/go-sqlite3"
)

func CreatePost(db *sql.DB, Title, Content, Date string, Sender int, Image, Like, Dislike string) error {
	// Insertion dans la base de données
	statement, err := db.Prepare("INSERT INTO Post (Title, Content, Date, Sender, Image, Like, Dislike) VALUES (?, ?, ?, ?, ?, ?, ?)")
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

// Récupère un post par son id
func SelectPost(db *sql.DB, PoID int) (database.Post, error) {
	var post database.Post
	err := db.QueryRow("SELECT id, Title, Content, Date, Sender, Image, Like, Dislike FROM Post WHERE id = ?", PoID).Scan(
		&post.PostID, &post.Title, &post.Content, &post.Date, &post.Sender, &post.Image, &post.Like, &post.Dislike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return post, errors.New("post non trouvé")
		}
		return post, err
	}
	return post, nil
}

// Function to get UserID based on Username
func GetUserIDByUsername(db *sql.DB, username string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT id FROM Users WHERE Pseudo = ?", username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
