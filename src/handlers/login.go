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
	if IsCookieExist(r) {
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
	// Redirect on home if there no POST method
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Retrieve the user from the database
	user, err := models.SelectUser(db, username)
	if err != nil || models.CheckPassword(user.Password, password) != nil {
		// Passer un message d'erreur au template
		data := map[string]interface{}{
			"ErrorMessage": "Nom d'utilisateur ou mot de passe incorrect",
		}
		tmpl := template.Must(template.ParseFiles("./web/template/login.html"))
		tmpl.Execute(w, data)
		return

	}

	// Create a session for the valid login
	SetSession(w, r, username)

	// Redirect to the dashboard after successful login
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Handle logout and clear session
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ClearSession(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
