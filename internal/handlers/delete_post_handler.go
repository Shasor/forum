package handlers

import (
	"fmt"
	"forum/internal/db"
	"net/http"
	"strconv"
)

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
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

	// Analyse les données du formulaire
	if err := r.ParseForm(); err != nil {
		fmt.Println("ParseForm error:", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Récupère l'ID du post à supprimer
	postIDStr := r.FormValue("id-post-to-delete")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		panic(err)
	}

	// Supprime le post en appelant la fonction appropriée
	if err := db.DeletePostByID(postID); err != nil {
		Resp.Msg = append(Resp.Msg, err.Error())
	}

	if Resp.Msg == nil {
		Resp.Msg = append(Resp.Msg, "Post successfully deleted!")
	}
	// Redirige après la suppression
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
