package db

import (
	"database/sql"
	"errors"
	"log"

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

func SelectUserByID(userID int) (User, error) {
	db := GetDB()
	defer db.Close()

	var user User
	err := db.QueryRow(`
	SELECT u.id, u.role, u.username, u.email, u.picture, u.password
	FROM users u
	WHERE id = ?`,
		userID).Scan(&user.ID, &user.Role, &user.Username, &user.Email, &user.Picture, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}

	user.Follows = GetUserFollows(user.ID)

	return user, nil
}

func SelectUserByUsername(username string) (User, error) {
	db := GetDB()
	defer db.Close()

	var user User
	err := db.QueryRow(`
	SELECT u.id, u.role, u.username, u.email, u.picture, u.password
	FROM users u
	WHERE username = ?`,
		username).Scan(&user.ID, &user.Role, &user.Username, &user.Email, &user.Picture, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}

	user.Follows = GetUserFollows(user.ID)

	return user, nil
}

func DeleteUserByUsername(username string) error {
	db := GetDB()
	defer db.Close()

	// Retrieve the user ID to delete
	var userID int
	err := db.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&userID)
	if err != nil {
		log.Printf("Error retrieving user ID: %v", err)
		return err
	}

	// Update posts by this user to set sender to 0
	_, err = db.Exec(`UPDATE posts SET sender = 0 WHERE sender = ?`, userID)
	if err != nil {
		log.Printf("Error updating posts sender to 0: %v", err)
		return err
	}

	// Now delete the user
	_, err = db.Exec(`DELETE FROM users WHERE id = ?`, userID)
	if err != nil {
		log.Printf("Error when deleting user: %v", err)
		return err
	}

	return nil
}

// UpdateUserProfile updates the user's email and picture in the database.
func UpdateUserProfile(userID int, email, picture string) error {
	db := GetDB()
	defer db.Close()

	// Prepare the SQL statement
	stmt, err := db.Prepare("UPDATE users SET email = ?, picture = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement with the provided values
	_, err = stmt.Exec(email, picture, userID)
	return err
}

func IsPasswordValid(providedPassword, storedHash string) bool {
	// Comparer le mot de passe fourni avec le hash stocké
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(providedPassword))

	// Si err est nil, cela signifie que les mots de passe correspondent
	return err == nil
}

func GetUserFollows(id int) []Category {
	db := GetDB()
	defer db.Close()

	query := `
        SELECT c.id, c.name
        FROM categories c
		JOIN follows f
		WHERE c.id = f.category AND f.user = ?;`

	rows, err := db.Query(query, id)
	if err != nil {
		// Gérer l'erreur, par exemple en la journalisant
		log.Printf("Erreur lors de l'exécution de la requête : %v", err)
		return nil
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			// Gérer l'erreur, par exemple en la journalisant
			log.Printf("Erreur lors du scan de la ligne : %v", err)
			continue
		}
		categories = append(categories, category)
	}

	return categories
}
