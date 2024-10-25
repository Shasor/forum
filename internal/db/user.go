package db

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(role, username, email, picture, password string) (*User, error) {
	// Ouvrir la connexion à la base de données
	db := GetDB()
	defer db.Close()

	// Démarrer une transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	// Préparer la requête d'insertion
	stmt, err := tx.Prepare("INSERT INTO users(role, username, email, picture, password) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Exécuter l'insertion
	_, err = stmt.Exec(role, username, email, picture, password)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("l'email ou le pseudo existe déjà")
	}

	// Commit de la transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	user, _ := SelectUserByUsername(username)
	return &user, nil
}

func SelectUserByUsername(username string) (User, error) {
	db := GetDB()
	defer db.Close()

	var user User
	err := db.QueryRow(`
	SELECT id, role, username, email, picture, password
	FROM users
	WHERE username = ?`,
		username).Scan(&user.ID, &user.Role, &user.Username, &user.Email, &user.Picture, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}

	return user, nil
}

func IsPasswordValid(providedPassword, storedHash string) bool {
	// Comparer le mot de passe fourni avec le hash stocké
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(providedPassword))

	// Si err est nil, cela signifie que les mots de passe correspondent
	return err == nil
}
