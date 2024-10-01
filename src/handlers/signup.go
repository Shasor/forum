package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"login/src/models"
)

// Affiche la page d'inscription
func SignupPage(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if IsLoggedIn(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/template/signup.html"))
		tmpl.Execute(w, nil)
	} else if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Créer un nouvel utilisateur
		err := models.CreateUser(db, username, password)
		if err != nil {
			http.Error(w, "Erreur lors de la création du compte", http.StatusInternalServerError)
			return
		}

		// Rediriger vers la page de connexion
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
