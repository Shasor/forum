package handlers

import (
	"forum/internal/db"
	"strconv"
	"net/http"
	"encoding/json"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID de l'utilisateur depuis l'URL
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Resp.Msg = append(Resp.Msg, "ID utilisateur invalide")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user, _ := db.SelectUserById(id)
	userPosts := db.GetPostFromUserById(id)

	userPostsReactions := db.FetchPostsReactions(id)
	userPosts = postsUnion(userPosts, userPostsReactions)

	db.SortPostsByDateDesc(userPosts)


	dataUser := map[string]interface{}{
		"userData": user,
		"userPosts": userPosts,
		//"userLikedPosts": userLikedPosts,
	}

	// Retourner les informations de l'utilisateur en format JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dataUser)
}

func postsUnion(posts1, posts2 []db.Post) []db.Post{

	for _, post := range posts2{
		posts1 = append(posts1, post)
	}

	return posts1

}