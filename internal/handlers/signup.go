package handlers

import (
	"forum/internal/db"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther) // not sure good code
	}

	var errors []string
	email := r.FormValue("email")
	username := r.FormValue("pseudo")
	password, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		errors = append(errors, err.Error())
	}
	picture, _ := db.ImageToBase64("static/assets/img/default_profile_picture.png")

	err = db.CreateUser("user", username, email, picture, string(password))
	if err != nil {
		errors = append(errors, err.Error())
	}

	tmpl, err := template.ParseFiles("web/pages/index.html", "web/templates/header.html", "web/templates/left-bar.html", "web/templates/posts.html", "web/templates/create-post.html", "web/templates/js.html")
	if err != nil {
		http.Error(w, "Internal Server Error (Error parsing templates)", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"err": errors,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error (Error executing template)", http.StatusInternalServerError)
		return
	}
}
