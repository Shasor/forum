package handlers

import (
	"database/sql"
	"html/template"
	"login/src/models"
	"net/http"
)

// Display the login page
func LoginPage(w http.ResponseWriter, r *http.Request) {
	// Redirect to dashboard if already logged in
	if IsLoggedIn(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// Serve the login page
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./web/template/login.html"))
		tmpl.Execute(w, nil)
	}
}

// Handle login form submission
func LoginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Retrieve the user from the database
		user, err := models.GetUserByUsername(db, username)
		if err != nil || models.CheckPassword(user.Password, password) != nil {
			http.Error(w, "Nom d'utilisateur ou mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		// Create a session for the valid login
		SetSession(w, r, username)

		// Redirect to the dashboard after successful login
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// Handle logout and clear session
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ClearSession(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
