package handlers

import (
	"fmt"
	"forum/internal/db"
	"net/http"
	"strconv"
	"time"
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
	date := fmt.Sprintf("%02d:%02d | %02d/%02d/%d", time.Now().Hour(), time.Now().Minute(), time.Now().Day(), time.Now().Month(), time.Now().Year())

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

	senderID := db.CheckReport(postID)
	if senderID != 0 {
		db.AddNotification("reportdone", date, 0, senderID, postID, 0)
		db.ReadNotification("report", senderID, 0, postID)
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
