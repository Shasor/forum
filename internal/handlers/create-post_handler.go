package handlers

import (
	"fmt"
	"forum/internal/db"
	"net/http"
	"strconv"
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

	sender, _ := strconv.Atoi(r.FormValue("sender_post"))
	category := r.FormValue("categorie_post")
	title := r.FormValue("title_post")
	content := r.FormValue("content_post")

	Resp = Response{}
	// Check if the form has a file for image_post
	var base64image string
	if file, header, err := r.FormFile("image_post"); err == nil {
		defer file.Close()
		// If header is not empty, encode the image
		if header.Size > 0 {
			base64image, err = ImageToBase64(file, header)
			if err != nil {
				Resp.Msg = append(Resp.Msg, err.Error())
			}
		}
	}

	date := fmt.Sprintf("%02d:%02d | %02d/%02d/%d", time.Now().Hour(), time.Now().Minute(), time.Now().Day(), time.Now().Month(), time.Now().Year())

	_ = db.CreatePost(sender, category, title, content, base64image, date)

	Resp.Msg = append(Resp.Msg, "Your post has been successfully sent!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
