package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitAllDB(db *sql.DB) {
	// Création de la table Users
	UserTable, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS Users 
		(id INTEGER PRIMARY KEY AUTOINCREMENT,
		Email TEXT UNIQUE NOT NULL,
		Pseudo TEXT UNIQUE NOT NULL,
		Password TEXT NOT NULL,
		Role TEXT NOT NULL,
		ProfilePicture TEXT,
		FollowID TEXT);
	`)

	// Création de la table Post avec la référence à Categories
	PostTable, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS Post 
		(id INTEGER PRIMARY KEY AUTOINCREMENT,
		CategoryID INTEGER,
		Title TEXT NOT NULL,
		Content TEXT NOT NULL,
		Date TEXT NOT NULL,
		Sender INTEGER,
		Image TEXT,
		Like TEXT,
		Dislike TEXT,
		FOREIGN KEY (Sender) REFERENCES Users(id),
		FOREIGN KEY (CategoryID) REFERENCES Categories(id));
	`)

	// Création de la table Commentaire
	CommentaireTable, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS Commentaire 
		(id INTEGER PRIMARY KEY AUTOINCREMENT,
		SenderID INTEGER,
		PostID INTEGER,
		Like TEXT,
		Dislike TEXT,
		Date TEXT NOT NULL,
		Content TEXT NOT NULL,
		FOREIGN KEY (SenderID) REFERENCES Users(id),
		FOREIGN KEY (PostID) REFERENCES Post(id));
	`)

	// Création de la table Categories
	CategorieTable, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS Categories 
		(id INTEGER PRIMARY KEY AUTOINCREMENT,
		Name TEXT NOT NULL);
	`)

	// Exécution des requêtes pour créer les tables
	UserTable.Exec()
	PostTable.Exec()
	CommentaireTable.Exec()
	CategorieTable.Exec()
}
