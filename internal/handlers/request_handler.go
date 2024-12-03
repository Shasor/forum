package handlers

import (
	"fmt"
	"forum/internal/db"
	"net/http"
	"time"
)

func RequestHandler(w http.ResponseWriter, r *http.Request) {
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

	user := GetUserFromCookie(w, r)

	err := db.AddNotification("request", date, user.ID, 0, 0, 0)
	if err != nil {
		Resp.Msg = append(Resp.Msg, err.Error())
	}

	if Resp.Msg == nil {
		Resp.Msg = append(Resp.Msg, "Request successfully sent to moderators!")
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
