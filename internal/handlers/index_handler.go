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
	data := map[string]interface{}{
		"resp":  Resp,
		"user":  GetUserFromCookie(w, r),
		"posts": db.FetchPosts(),
	}
	Parse(w, data)
}
