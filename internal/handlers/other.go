package handlers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"forum/internal/db"
	"image"
	_ "image/gif" // Import pour le décodage GIF
	"image/jpeg"
	_ "image/png" // Import pour le décodage PNG
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

func GetFormGET(w http.ResponseWriter, r *http.Request) Get {
	var get Get
	if categoryIDStr := r.URL.Query().Get("catID"); categoryIDStr != "" {
		var err error
		get.CategoryID, err = strconv.Atoi(categoryIDStr)
		catExist := db.CategoryExist(db.GetCategoryNameByID(get.CategoryID))
		if !catExist {
			Resp.Msg = append(Resp.Msg, "The category you're asking for doesn't exist!")
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		get.Type = "category"
		if err != nil {
			get.CategoryID = 0
		}
	} else {
		get.CategoryID = 0
	}
	if postStr := r.URL.Query().Get("postID"); postStr != "" && IsCookieValid(w, r) {
		var err error
		get.PostID, err = strconv.Atoi(postStr)
		postExist := db.PostExist(get.PostID)
		if !postExist {
			Resp.Msg = append(Resp.Msg, "The post you're asking for doesn't exist!")
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		get.Type = "post"
		if err != nil {
			get.PostID = 0
		}
	} else if !IsCookieValid(w, r) && postStr != "" {
		get.PostID = 0
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		// fmt.Println("Redirect test", postStr)
		get.PostID = 0
	}
	return get
}

func ImageToBase64(file multipart.File, header *multipart.FileHeader, is_pfp bool) (string, error) {
	// Reset the file pointer to the beginning of the file
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("error resetting file pointer: %v", err)
	}
	// Check the file extension to see if it's a GIF
	ext := filepath.Ext(header.Filename)
	if ext != ".gif" && ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		return "", fmt.Errorf("invalid file extension: %s", ext)
	}
	// Check that the image does not exceed 20mb
	if header.Size > 20000000 {
		return "", errors.New("the image is over 20mb")
	}
	if ext == ".gif" {
		// For GIFs, just read the entire file and base64 encode it
		data, err := io.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("error reading GIF file: %v", err)
		}
		return base64.StdEncoding.EncodeToString(data), nil
	}

	// For non-GIF images (JPG, PNG), proceed with decoding and resizing
	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("error decoding image: %v", err)
	}
	var resizedImg image.Image
	// Resize the image (optional, depends on your requirements)
	if is_pfp {
		resizedImg = resizeImage(img, 500, 500)
	} else {
		resizedImg = resizeImage(img, 800, 800)
	}

	// Encode resized image to JPEG format
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resizedImg, nil)
	if err != nil {
		return "", fmt.Errorf("error encoding image to JPEG: %v", err)
	}

	// Base64-encode the JPEG image
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// Resize img to a specified maxWidth and maxHeight
func resizeImage(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Calculate the resize ratio
	ratioW := float64(maxWidth) / float64(width)
	ratioH := float64(maxHeight) / float64(height)
	ratio := math.Min(ratioW, ratioH)

	newWidth := int(float64(width) * ratio)
	newHeight := int(float64(height) * ratio)

	// Create a new resized image
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Resize the image manually
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x) / ratio)
			srcY := int(float64(y) / ratio)
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}
	return dst
}

func normalizeSpaces(s string) string {
	r := strings.Fields(s)
	return strings.Join(r, " ")
}

func capitalize(s string) string {
	// Diviser la chaîne en mots
	words := strings.Fields(s)

	// Parcourir chaque mot
	for i, word := range words {
		// Convertir le premier caractère en majuscule et le reste en minuscule
		runes := []rune(word)
		for j := range runes {
			if j == 0 {
				runes[j] = unicode.ToUpper(runes[j])
			} else {
				runes[j] = unicode.ToLower(runes[j])
			}
		}
		words[i] = string(runes)
	}

	// Rejoindre les mots en une seule chaîne
	return strings.Join(words, " ")
}

func GetFileFromURL(url string) (string, error) {
	// Effectuer une requête GET vers le lien
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la requête GET : %v", err)
	}
	defer resp.Body.Close()

	// Vérifier le code de statut HTTP
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("échec de la requête : statut HTTP %d", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", err
	}
	resizedImg := resizeImage(img, 500, 500)

	// Encode resized image to JPEG format
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resizedImg, nil)
	if err != nil {
		return "", fmt.Errorf("error encoding image to JPEG: %v", err)
	}

	// Base64-encode the JPEG image
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
