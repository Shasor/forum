package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"login/src/database"
	"login/src/models"
	"net/http"
)

// EditProfilePage handler for editing user profile
func EditProfilePage(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	// Handle form submission
	if r.Method == http.MethodPost {
		newEmail := r.FormValue("email")
		var newProfilePicture string

		if file, _, err := r.FormFile("profile_picture"); err == nil {
			newProfilePicture, err = base64Image(file) // Convert the uploaded file to base64
			if err != nil {
				log.Println("Error processing profile picture:", err) // Log the error
				http.Error(w, "Error processing profile picture", http.StatusBadRequest)
				return
			}
		} else {
			log.Println("No profile picture uploaded, keeping old one.") // Log to confirm no file was uploaded
			newProfilePicture = user.ProfilePicture
		}

		// Update user information
		if newEmail == "" {
			newEmail = user.Email // Keep the old email if none is provided
		}

		user.Email = newEmail
		user.ProfilePicture = newProfilePicture

		log.Printf("Attempting to update user: Email=%s, ProfilePictureLength=%d, UserID=%d\n", user.Email, len(user.ProfilePicture), user.UserID)

		if err := models.UpdateUser(db, user); err != nil {
			log.Println("Error updating user:", err) // Log the error for debugging
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}

		log.Printf("After updating user : Email=%s, ProfilePictureLength=%d, UserID=%d\n", user.Email, len(user.ProfilePicture), user.UserID)

		// Redirect to the profile page after a successful update
		http.Redirect(w, r, "/profile/"+user.Pseudo, http.StatusSeeOther)
		return
	}

	// Render edit profile page with user data
	data := struct {
		User database.User
	}{User: user}

	tmpl := template.Must(template.ParseFiles("./web/template/edit-profile.html"))
	tmpl.Execute(w, data)
}
