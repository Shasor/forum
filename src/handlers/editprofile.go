package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"login/src/database"
	"login/src/models"
	"mime/multipart"
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

	var errorMessage string

	// Handle form submission
	if r.Method == http.MethodPost {
		log.Printf("Attempting to update user: Email=%s, ProfilePictureLength=%d, UserID=%d\n", user.Email, len(user.ProfilePicture), user.UserID)
		newEmail := r.FormValue("email")
		var newProfilePicture string

		if file, _, err := r.FormFile("profile_picture"); err == nil {
			// Check file size
			if err := checkFileSize(file); err != nil {
				errorMessage = err.Error() // Set the error message if the file is too large
			} else {
				newProfilePicture, err = base64Image(file) // Convert the uploaded file to base64
				if err != nil {
					log.Println("Error processing profile picture:", err) // Log the error
					errorMessage = "Error processing profile picture."
				}
			}
		} else {
			log.Println("No profile picture uploaded, keeping old one.") // Log to confirm no file was uploaded
			newProfilePicture = user.ProfilePicture
		}

		// If there's an error, don't proceed with the update
		if errorMessage != "" {
			// Render edit profile page with user data and error message
			data := struct {
				User         database.User
				ErrorMessage string
			}{User: user, ErrorMessage: errorMessage}

			tmpl := template.Must(template.ParseFiles("./web/template/edit-profile.html"))
			tmpl.Execute(w, data)
			return
		}

		// Update user information
		if newEmail == "" {
			newEmail = user.Email // Keep the old email if none is provided
		}

		user.Email = newEmail
		user.ProfilePicture = newProfilePicture

		if err := models.UpdateUser(db, user); err != nil {
			log.Println("Error updating user:", err) // Log the error for debugging
			errorMessage = "Error updating user."
		} else {
			// Redirect to the profile page after a successful update
			http.Redirect(w, r, "/profile/"+user.Pseudo, http.StatusSeeOther)
			return
		}
	}

	// Render edit profile page with user data
	data := struct {
		User         database.User
		ErrorMessage string
	}{User: user, ErrorMessage: errorMessage}

	tmpl := template.Must(template.ParseFiles("./web/template/edit-profile.html"))
	tmpl.Execute(w, data)
}

// Function to check file size
func checkFileSize(file multipart.File) error {
	const maxFileSize = 8 * 1024 * 1024 // 8 MB
	buf := make([]byte, maxFileSize)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}

	if n == maxFileSize && err == nil {
		return fmt.Errorf("Uploaded file exceeds the maximum size limit of 8MB.")
	}
	return nil
}
