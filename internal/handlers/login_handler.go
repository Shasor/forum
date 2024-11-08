package handlers

import (
	"forum/internal/db"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || IsCookieValid(w, r) {
		Resp = Response{Msg: []string{
			map[bool]string{
				true:  "You are already connected",
				false: "Method not Allowed",
			}[IsCookieValid(w, r)],
		}}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	Resp = Response{}
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := db.SelectUserByUsername(username)
	if username != "" && password != "" {
		if err == nil {
			if db.IsPasswordValid(password, user.Password) {
				SetSession(w, username)
			} else {
				Resp.Msg = append(Resp.Msg, "Wrong password!")
				Resp.Action = "GetLogin();"
			}
		} else {
			Resp.Msg = append(Resp.Msg, err.Error())
			Resp.Action = "GetLogin();"
		}
	} else {
		Resp.Msg = append(Resp.Msg, "All fields are required!")
		Resp.Action = "GetLogin();"
	}

	if Resp.Msg == nil {
		Resp.Msg = append(Resp.Msg, "You've successfully logged in!")
	} else {
		Resp.Form.Username = username
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
