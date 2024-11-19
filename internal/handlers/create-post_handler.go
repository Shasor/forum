package handlers

import (
	"fmt"
	"forum/internal/db"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST and if the user is connected (cookie is valid)
	if r.Method != http.MethodPost || !IsCookieValid(w, r) {
		Resp = Response{Msg: []string{
			map[bool]string{
				true:  "Method not Allowed",
				false: "You are not connected",
			}[r.Method != http.MethodPost],
		}}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Initialize the response object
	Resp = Response{}

	// Retrieve form values
	sender, _ := strconv.Atoi(r.FormValue("sender_post"))
	category := capitalize(r.FormValue("categorie_post"))
	title := normalizeSpaces(r.FormValue("title_post"))
	content := normalizeSpaces(r.FormValue("content_post"))

	// Check if the form has an image (if it does, convert it to base64)
	var base64image string
	if file, header, err := r.FormFile("image_post"); err == nil {
		defer file.Close()
		// If header is not empty, encode the image
		if header.Size > 0 {
			base64image, err = ImageToBase64(file, header, false)
			if err != nil {
				Resp.Msg = append(Resp.Msg, err.Error())
			}
		}
	}

	// Validate form data
	if category == "" {
		Resp.Msg = append(Resp.Msg, "All fields must be completed!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else if content == "" && base64image == "" {
		// Both content and image are missing
		Resp.Msg = append(Resp.Msg, "Must fill every box")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else if len(category) >= 26 {
		Resp.Msg = append(Resp.Msg, "The category must not exceed 25 characters!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Format the date in the desired format
	date := fmt.Sprintf("%02d:%02d | %02d/%02d/%d", time.Now().Hour(), time.Now().Minute(), time.Now().Day(), time.Now().Month(), time.Now().Year())

	// Retrieve the parent_id from the form (if this is a comment, parent_id will be set)
	var parentID *int // Pointer to allow for nil value (original posts have nil parentID)
	parentIDParam := r.FormValue("parent_id")
	if parentIDParam != "" {
		// If parent_id is provided, convert it to an integer
		parentIDValue, err := strconv.Atoi(parentIDParam)
		if err != nil {
			log.Println("Invalid parent_id:", err)
			Resp.Msg = append(Resp.Msg, "Invalid parent_id")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		parentID = &parentIDValue // Set parent_id for a comment
	} else {
		// If parent_id is not provided, it's an original post, set parent_id to nil (or 0 inside the db)
		parentID = nil
	}

	category = PurgeHashAndSpace(category)

	categories := strings.Split(category, "#")
	postAlreadyCreated := false
	for _, cat := range categories {
		if !db.CategoryExist(capitalize(cat)) {
			db.CreateCategory(capitalize(cat))
		}
		if !postAlreadyCreated {
			_ = db.CreatePost(sender, cat, title, content, base64image, date, parentID)
			postAlreadyCreated = true
		}
		category_cat, _ := db.SelectCategoryByName(capitalize(cat))
		postID, _ := db.GetLastPostIDByUserID(sender)
		db.LinkPostToCategory(postID, category_cat)
	}

	Resp.Msg = append(Resp.Msg, "Your post has been successfully sent!")
	if parentID != nil {
		id := strconv.Itoa(*parentID)
		http.Redirect(w, r, "/?postID="+id, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
