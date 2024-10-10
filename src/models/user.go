package models

import (
	"database/sql"
	"errors"
	"log"
	"login/src/database"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Récupère un utilisateur par son nom d'utilisateur
func SelectUser(db *sql.DB, Pseudo string) (database.User, error) {
	var user database.User
	err := db.QueryRow("SELECT id, Email, Pseudo, Password, Role, ProfilePicture, FollowID FROM Users WHERE Pseudo = ?", Pseudo).Scan(
		&user.UserID, &user.Email, &user.Pseudo, &user.Password, &user.Role, &user.ProfilePicture, &user.FollowID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("utilisateur non trouvé")
		}
		return user, err
	}
	return user, nil
}

func IsUserPresent(db *sql.DB, username string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Users WHERE Pseudo = ?", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Ajoute un nouvel utilisateur dans la base de données avec un mot de passe hashé
func CreateUser(db *sql.DB, Email, Pseudo, Password, Role, ProfilePicture, FollowID string) error {
	// Hasher le mot de passe avant de l'enregistrer
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insertion dans la base de données
	statement, err := db.Prepare("INSERT INTO Users (Email, Pseudo, Password, Role, ProfilePicture, FollowID) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(Email, Pseudo, hashedPassword, Role, ProfilePicture, FollowID)
	return err
}

// Supprimer un utilisateur par son nom d'utilisateur
func DeleteUserByUsername(db *sql.DB, Pseudo string) error {
	// Requête SQL pour supprimer l'utilisateur
	stmt, err := db.Prepare("DELETE FROM Users WHERE Pseudo = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(Pseudo)
	if err != nil {
		log.Println("Erreur lors de la suppression de l'utilisateur :", err)
		return err
	}

	return nil
}

func UserExist(db *sql.DB, pseudo string) bool {
	var exist bool
	err := db.QueryRow("SELECT EXISTS( SELECT 1 FROM Users WHERE Pseudo = ?)", pseudo).Scan(&exist)
	if err != nil {
		return false
	}
	return exist

}

// Function to get UserID based on Username
func GetUsernameByID(db *sql.DB, id int) (string, error) {
	var username = ""
	err := db.QueryRow("SELECT pseudo FROM Users WHERE id = ?", id).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

// Vérifie si le mot de passe fourni correspond au hash enregistré dans la base de données
func CheckPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func UpdateUser(db *sql.DB, user database.User) error {

	query := `UPDATE users SET Email = ?, ProfilePicture = ? WHERE id = ?`
	_, err := db.Exec(query, user.Email, user.ProfilePicture, user.UserID) // Ensure user.UserID is correctly passed
	if err != nil {
		log.Println("Error executing UPDATE query:", err) // Log any SQL execution error
	}
	return err
}

/* Below are functions for FollowID


 */

func UpdateFollowedID(db *sql.DB, UserID, FollowID string) error {
	statement, err := db.Prepare("UPDATE Users SET FollowID = ? Where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(FollowID, UserID)
	if err != nil {
		log.Println("Erreur lors de la mise à jour des abonnements :", err)
		return err
	}
	return nil
}

func AddFollowID(FollowID string, newFollowID int) string {
	newString := strconv.Itoa(newFollowID)
	return FollowID + "-" + newString
}

func SplitFollowedID(categoriesID string) []string {
	tab := strings.Split(categoriesID, "-")
	return tab
}

func RemoveFollowedID(cateID []string, id int) string {
	idString := strconv.Itoa(id)
	categoriesID := ""
	for indice, value := range cateID {
		if value == idString {
			continue
		}
		categoriesID += value
		if indice <= len(cateID)-1 {
			categoriesID += "-"
		}
	}
	return categoriesID
}
