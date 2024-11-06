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
	// Get the user from the cookie, which may be nil
	user := GetUserFromCookie(w, r)
	if user != nil {
		likedposts = db.FetchPostsLiked(user.ID)
	}
	categoryID, postID := GetFormGet(r)
	data := map[string]interface{}{
		"resp":       Resp,
		"user":       user, // This will be nil if no user is logged in
		"posts":      db.FetchPosts(categoryID, postID),
		"likedposts": likedposts,
	}

	Parse(w, data)
}