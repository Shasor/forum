package handlers

import (
	"database/sql"
	"log"
	"login/src/database"
	"login/src/models" // Update with the correct path to your models
	"net/http"
)

// DeleteAccountHandler handles user account deletion
func DeleteAccountHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if !IsCookieExist(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Verify method
	if r.Method != "POST" {
		log.Println("Invalid method:", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, err := GetSessionUsername(r)
	if err != nil {
		log.Println("No user logged in.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Log username
	log.Println("Deleting user:", username)

	// Delete user from the database
	err = models.DeleteUserByUsername(db, username)
	if err != nil {
		log.Println("Error deleting account:", err)
		http.Error(w, "Unable to delete account", http.StatusInternalServerError)
		return
	}

	// Log success
	log.Println("Account successfully deleted:", username)

	// Clear user session
	ClearSession(w, r)

	// Display confirmation page
	RenderTemplate(w, "account-deleted.html", nil)
}

// UserProfilePage displays the profile of a specific user by their username
func UserProfilePage(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	username, err := GetSessionUsername(r) // Get the username from the session
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := models.SelectUser(db, username) // Retrieve the user
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Prepare the data for the template
	data := struct {
		User database.User // User information
	}{
		User: user,
	}

	RenderTemplate(w, "profile.html", data) // Render the profile template
}

// UserPostsPage displays the posts of a specific user
func UserPostsPage(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	username, err := GetSessionUsername(r) // Get the username from the session
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := models.SelectUser(db, username) // Retrieve the user
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Fetch posts for the current user (by their ID)
	posts, err := models.FetchPostsBySender(db, user.UserID) // Assuming user.ID is the Sender ID
	if err != nil {
		http.Error(w, "Unable to fetch posts", http.StatusInternalServerError)
		return
	}

	// Prepare the data for the template
	data := struct {
		User  database.User   // User information
		Posts []database.Post // User's posts
	}{
		User:  user,
		Posts: posts,
	}

	RenderTemplate(w, "posts.html", data) // Render the posts template
}
