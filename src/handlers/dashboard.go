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
	"login/src/database"
	"login/src/models"
	"math"
	"mime/multipart"
	"net/http"
	"time"
)

// Page du dashboard avec le nom d'utilisateur
func DashboardPage(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if !IsCookieExist(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username, err := GetSessionUsername(r)
	if err != nil {
		if err == ErrInvalidCookie {
			ClearSession(w, r) // If the Cookie hash is incorrect, delete the cookie
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther) // If no session, redirect to login
		return
	}

	// Fetch users from the database (limit to 5)
	users, err := FetchUsers(db, 5, username)
	if err != nil {
		log.Println("Error fetching users:", err)
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("./web/template/dashboard.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Username": username,
		"Users":    users,
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

	// image management
	file, _, err := r.FormFile("image_post")
	ErrorsExit(err)
	defer file.Close()
	base64image, err := base64Image(file)
	ErrorsExit(err)

	date := time.Now().Format(layout)

	senderObject, err := models.SelectUser(db, sender)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Call your model function to insert the post
	err = models.CreatePost(db, title, content, date, senderObject.UserID, base64image, "0", "0")
	if err != nil {
		// Log the error for debugging
		log.Println("Error creating post:", err)
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Convert image file into string in base64
func base64Image(file multipart.File) (string, error) {
	// Décodez l'image
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// Redimensionnez l'image
	resizedImg := resizeImage(img)

	// Encodez l'image redimensionnée en JPEG
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resizedImg, nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// Resize img to 1920x1080
func resizeImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	maxWidth := 1920
	height := bounds.Dy()
	maxHeight := 1080

	// Calculez le ratio pour redimensionner
	ratioW := float64(maxWidth) / float64(width)
	ratioH := float64(maxHeight) / float64(height)
	ratio := math.Min(ratioW, ratioH)

	if ratio >= 1 {
		return img
	}

	newWidth := int(float64(width) * ratio)
	newHeight := int(float64(height) * ratio)

	// Créez une nouvelle image redimensionnée
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Redimensionnez l'image manuellement
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x) / ratio)
			srcY := int(float64(y) / ratio)
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}

	return dst
}

// FetchUsers fetches a limited number of users from the database, excluding the current user's account
func FetchUsers(db *sql.DB, limit int, username string) ([]database.User, error) {
	// Use a parameterized query to prevent SQL injection
	query := "SELECT Pseudo, ProfilePicture FROM Users WHERE Pseudo != ? LIMIT ?"
	rows, err := db.Query(query, username, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []database.User
	for rows.Next() {
		var user database.User
		err := rows.Scan(&user.Pseudo, &user.ProfilePicture)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
