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
	"math"
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
	defer file.Close()
	base64image, err := base64Image(file)
	ErrorsHandler(err)
	// finished here last time

	date := time.Now().Format(layout)

	senderObject, err := models.SelectUser(db, sender)
	if err != nil {
		fmt.Println("Error SELECTING user in CreatePostHandler")
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

func ErrorsHandler(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
