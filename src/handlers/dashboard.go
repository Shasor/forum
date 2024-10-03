package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"login/src/models"
	"net/http"
	"time"
)

// Page du dashboard avec le nom d'utilisateur
func DashboardPage(w http.ResponseWriter, r *http.Request) {
	if !IsLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	username, exists := GetSessionUsername(r)
	if !exists {
		http.Redirect(w, r, "/", http.StatusSeeOther) // If no session, redirect to login
		return
	}

	tmpl := template.Must(template.ParseFiles("./web/template/dashboard.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Username": username,
	})
}

// CreatePostHandler handles the post creation
func CreatePostHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	layout := "2006-01-02 15:04:05"

	sender := r.FormValue("sender_post")
	title := r.FormValue("title_post")
	content := r.FormValue("content_post")
	image := r.FormValue("image_post")
	date := time.Now().Format(layout)

	senderObject, err := models.SelectUser(db, sender)
	if err != nil {
		fmt.Println("Error SELECTING user in CreatePostHandler")
		return
	}

	// Call your model function to insert the post
	err = models.CreatePost(db, title, content, date, senderObject.UserID, image, "0", "0")
	if err != nil {
		// Log the error for debugging
		log.Println("Error creating post:", err)
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
