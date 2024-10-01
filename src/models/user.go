package models

import (
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Structure représentant un utilisateur
type User struct {
	Username string
	Password string
}

// Récupère un utilisateur par son nom d'utilisateur
func GetUserByUsername(db *sql.DB, username string) (User, error) {
	var user User
	err := db.QueryRow("SELECT username, password FROM users WHERE username = ?", username).Scan(&user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("utilisateur non trouvé")
		}
		return user, err
	}
	return user, nil
}

// Ajoute un nouvel utilisateur dans la base de données avec un mot de passe hashé
func CreateUser(db *sql.DB, username, password string) error {
	// Hasher le mot de passe avant de l'enregistrer
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insertion dans la base de données
	statement, err := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(username, hashedPassword)
	return err
}

// Supprimer un utilisateur par son nom d'utilisateur
func DeleteUserByUsername(db *sql.DB, username string) error {
	// Requête SQL pour supprimer l'utilisateur
	stmt, err := db.Prepare("DELETE FROM users WHERE username = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username)
	if err != nil {
		log.Println("Erreur lors de la suppression de l'utilisateur :", err)
		return err
	}

	return nil
}

// Vérifie si le mot de passe fourni correspond au hash enregistré dans la base de données
func CheckPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
