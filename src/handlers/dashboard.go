package handlers

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"log"
	"login/src/models"
	"mime/multipart"
	"net/http"
	"os"
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
	file, _, err := r.FormFile("image_post")
	ErrorsHandler(err)
	encodedString, err := compressAndEncodeImage(file)
	ErrorsHandler(err)
	fmt.Println(len(encodedString))
	// finished here last time

	date := time.Now().Format(layout)

	senderObject, err := models.SelectUser(db, sender)
	if err != nil {
		fmt.Println("Error SELECTING user in CreatePostHandler")
		return
	}

	// Call your model function to insert the post
	err = models.CreatePost(db, title, content, date, senderObject.UserID, encodedString, "0", "0")
	if err != nil {
		// Log the error for debugging
		log.Println("Error creating post:", err)
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func compressAndEncodeImage(file multipart.File) (string, error) {
	// Décode l'image
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// Prépare un buffer pour l'image compressée
	buf := new(bytes.Buffer)

	// Compresse l'image en JPEG avec une qualité réduite
	err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 50})
	if err != nil {
		return "", err
	}

	// Encode en base64
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// func compressString(s string) (string, error) {
// 	var b bytes.Buffer
// 	gz := gzip.NewWriter(&b)
// 	if _, err := gz.Write([]byte(s)); err != nil {
// 		return "", err
// 	}
// 	if err := gz.Close(); err != nil {
// 		return "", err
// 	}
// 	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
// }

func ErrorsHandler(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
