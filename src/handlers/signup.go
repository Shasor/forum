package handlers

import (
	"database/sql"
	"encoding/base64"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"login/src/models"
)

// Affiche la page d'inscription
func SignupPage(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if IsCookieExist(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/template/signup.html"))
		tmpl.Execute(w, nil)
	} else if r.Method == "POST" {
		Email := r.FormValue("email")
		Pseudo := r.FormValue("pseudo")
		Password := r.FormValue("password")

		// Get the default profile picture as a base64 string
		defaultImage, err := GetDefaultProfilePicture()
		if err != nil {
			// Handle the error appropriately
			defaultImage = "" // Fallback to an empty string or a different default
		}

		// Créer un nouvel utilisateur
		err = models.CreateUser(db, Email, Pseudo, Password, "User", defaultImage, "")
		if err != nil {
			// Passer un message d'erreur au template
			log.Println("Error when creating account:, ", err)
			data := map[string]interface{}{
				"ErrorMessage": "Erreur lors de la création du compte.",
			}
			tmpl := template.Must(template.ParseFiles("./web/template/signup.html"))
			tmpl.Execute(w, data)
			return
		}

		// Rediriger vers la page de connexion
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// GetDefaultProfilePicture returns the default profile picture as a base64 encoded string.
func GetDefaultProfilePicture() (string, error) {
	imgPath := "./web/img/default_profile_picture.jpg"
	imgData, err := ioutil.ReadFile(imgPath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(imgData), nil
}
