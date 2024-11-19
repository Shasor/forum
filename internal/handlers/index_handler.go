package handlers

import (
	"fmt"
	"forum/internal/db"
	"net/http"
	"strconv"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Erreur dans la requête", http.StatusBadRequest)
			return
		}

		// Récupère la valeur de l'input `search_bar`
		searchValue := r.FormValue("search_bar")

		//--TO DO-- : Recharger la page avec la catégorie si cela correspond, ou par l'ensemble des catégories contenant la recherche "ima" affiche les catégories "image", "Fatima", ...
		fmt.Println("Vous avez recherché : ", searchValue)
		categorySearched, err := db.SelectCategoryByName(searchValue)
		categoryID := strconv.Itoa(categorySearched.ID)
		if err != nil{
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		http.Redirect(w, r, "/?catID="+categoryID, http.StatusSeeOther)
	}

	if Resp.Broadcasted {
		Resp = Response{}
	} else {
		Resp.Broadcasted = true
	}
	var likedposts, followedposts []db.Post
	user := GetUserFromCookie(w, r)
	if user != nil {
		likedposts = db.FetchPostsLiked(user.ID)
		followedposts = db.FetchFollowPosts(user.ID)
	}
	get := GetFormGET(w, r)
	data := map[string]interface{}{
		"resp":          Resp,
		"user":          user, // This will be nil if no user is logged in
		"posts":         db.FetchPosts(),
		"followedposts": followedposts,
		"categories":    db.FetchCategories(),
		"likedposts":    likedposts,
		"GET":           get,
		"comments":      db.FetchComments(),
	}
	Parse(w, data)

}
