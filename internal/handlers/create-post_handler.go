package handlers

import (
	"fmt"
	"forum/internal/db"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
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
	Resp = Response{}

	sender, _ := strconv.Atoi(r.FormValue("sender_post"))
	category := capitalize(r.FormValue("categorie_post"))
	title := normalizeSpaces(r.FormValue("title_post"))
	content := normalizeSpaces(r.FormValue("content_post"))

	if category == "" || title == "" || content == "" {
		Resp.Msg = append(Resp.Msg, "All fields must be completed!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else if len(category) >= 26 {
		Resp.Msg = append(Resp.Msg, "The category must not exceed 25 characters!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Check if the form has a file for image_post
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

	date := fmt.Sprintf("%02d:%02d | %02d/%02d/%d", time.Now().Hour(), time.Now().Minute(), time.Now().Day(), time.Now().Month(), time.Now().Year())

	categories := strings.Split(category, "#")
	
	postAlreadyCreated := false
	for _, cat := range categories{
		if !postAlreadyCreated{
			_ = db.CreatePost(sender, cat, title, content, base64image, date)
			postAlreadyCreated = true
		}

		if !db.CategoryExist(capitalize(cat)){
			db.CreateCategory(capitalize(cat))
		}

		category_cat, _ := db.SelectCategoryByName(capitalize(cat))
		postID,_ := db.GetLastPostIDByUserID(sender)
		db.LinkPostToCategory(postID,category_cat )
	}
	

	Resp.Msg = append(Resp.Msg, "Your post has been successfully sent!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
