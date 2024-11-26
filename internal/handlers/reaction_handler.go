package handlers

import (
	"encoding/json"
	"forum/internal/db"
	"log"
	"net/http"
)

func ReactToPost(w http.ResponseWriter, r *http.Request) {
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

	var reqBody ReactionRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		panic(err)
	}

	postID := reqBody.PostID
	reaction := reqBody.Reaction

	user := GetUserFromCookie(w, r)
	if user == nil { // temporary
		return
	}

	err = db.UpdatePostReaction(user.ID, postID, reaction)
	if err != nil {
		log.Println("Error saving reaction:", err)
		panic(err)
	}

	// Fetch updated counts
	likes, dislikes, err := db.GetPostReactions(postID)
	if err != nil {
		log.Println("Error fetching updated reactions:", err)
		panic(err)
	}

	// Send updated counts back as JSON response
	response := map[string]int{
		"likes":    likes,
		"dislikes": dislikes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
