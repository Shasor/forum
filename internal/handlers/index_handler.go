package handlers

import (
	"forum/internal/db"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if Resp.Broadcasted {
		Resp = Response{}
	} else {
		Resp.Broadcasted = true
	}
	var likedposts []db.Post
	user := GetUserFromCookie(w, r)
	if user != nil {
		likedposts = db.FetchPostsLiked(user.ID)
	}
	get := GetFormGET(r)
	data := map[string]interface{}{
		"resp":       Resp,
		"user":       user, // This will be nil if no user is logged in
		"posts":      db.FetchPosts(),
		"categories": db.FetchCategories(),
		"likedposts": likedposts,
		"GET":        get,
	}
	Parse(w, data)

	if r.Method == "POST" {
		print(r.FormValue("SearchForm"))
	}

}
