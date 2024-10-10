package handlers

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
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

	username, err := GetSessionUsername(db, r)
	if err != nil {
		if err == ErrInvalidCookie {
			ClearSession(w, r) // If the Cookie hash is incorrect, delete the cookie
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther) // If no session, redirect to login
		return
	}

	// Fetch all categories from the database
	categories, err := models.FetchCategories(db)
	if err != nil {
		log.Println("Error fetching Categories:", err)
		http.Error(w, "Error fetching Categories", http.StatusInternalServerError)
		return
	}

	// Fetch all posts from the database
	posts, err := models.FetchPosts(db)
	if err != nil {
		log.Println("Error fetching Categories:", err)
		http.Error(w, "Error fetching Categories", http.StatusInternalServerError)
		return
	}
	// sends only requested posts
	var requestedPosts []database.Post
	if cat := r.URL.Query().Get("cat"); r.Method == http.MethodGet && cat != "" {
		catID, err := models.GetCategorieIDByName(db, cat)
		if err == nil {
			for _, post := range posts {
				if post.CategorieID == catID {
					senderUsername, _ := models.GetUsernameByID(db, post.SenderID)
					if senderUsername == "" {
						senderUsername = "Deleted User"
					}
					post.SenderUsername = senderUsername
					requestedPosts = append(requestedPosts, post)
				}
			}
		}
	}

	tmpl := template.Must(template.ParseFiles("./web/template/dashboard.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Username":   username,
		"Categories": categories,
		"Posts":      requestedPosts,
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
	categorie := r.FormValue("categorie_post")
	title := r.FormValue("title_post")
	content := r.FormValue("content_post")

	// Initialize base64image to an empty string
	var base64image string

	// Check if the form has a file for image_post
	if file, header, err := r.FormFile("image_post"); err == nil {
		defer file.Close() // Ensure file is closed after processing
		// If header is not empty, encode the image
		if header.Size > 0 {
			base64image, err = base64Image(file)
			ErrorsExit(err)
		}
	} else if err != http.ErrMissingFile {
		// Handle the error if it's not about a missing file
		ErrorsExit(err)
	}

	date := time.Now().Format(layout)

	senderObject, err := models.SelectUser(db, sender)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Call your model function to insert the post
	err = models.CreatePost(db, categorie, title, content, date, senderObject.UserID, base64image, "0", "0")
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
