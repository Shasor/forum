package handlers

import (
	"forum/internal/db"
	"net/http"
	"strings"
)

// EditProfileHandler handles updating the user profile.
func EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request is a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
		return
	}

	// Get the logged-in user from the session cookie
	user := GetUserFromCookie(w, r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse the form to retrieve new data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Retrieve values from the form
	email := strings.TrimSpace(r.FormValue("email"))
	pictureFile, header, err := r.FormFile("profile_picture")

	// Check for fields to skip if empty
	updatedEmail := user.Email     // Default to the current email if empty
	updatedPicture := user.Picture // Default to the current picture if not changed

	if email != "" {
		updatedEmail = email
	}

	// Update picture if a new file is uploaded
	if pictureFile != nil {
		defer pictureFile.Close()
		encodedPicture, err := ImageToBase64(pictureFile, header, true)
		if err != nil {
			http.Error(w, "Error encoding image", http.StatusInternalServerError)
			return
		}
		updatedPicture = encodedPicture
	}

	// Execute the database update
	err = db.UpdateUserProfile(user.ID, updatedEmail, updatedPicture)
	if err != nil {
		http.Error(w, "Error updating profile", http.StatusInternalServerError)
		return
	}

	// Redirect to the profile page or another appropriate page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
