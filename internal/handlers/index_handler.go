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

	// Get the user from the cookie, which may be nil
	user := GetUserFromCookie(w, r)

	categoryID, postID := GetFormGet(r)
	data := map[string]interface{}{
		"resp":       Resp,
		"user":       user, // This will be nil if no user is logged in
		"posts":      db.FetchPosts(categoryID, postID),
		"likedposts": db.FetchPostsLiked(user.ID),
	}

	Parse(w, data)
}
