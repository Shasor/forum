package handlers

import (
	"forum/internal/db"
	"net/http"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
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
	user := GetUserFromCookie(w, r)

	Resp = Response{Msg: []string{"Your account has successfully been deleted"}}
	ClearSession(w, r, "session_token")
	db.DeleteUserByUsername(user.Username)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
