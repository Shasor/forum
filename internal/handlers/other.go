package handlers

import (
	"forum/internal/db"
	"net/http"
	"strconv"
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
