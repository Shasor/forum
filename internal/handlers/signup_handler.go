package handlers

import (
	"forum/internal/db"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || IsCookieValid(w, r) {
		Resp = Response{Msg: []string{
			map[bool]string{
				true:  "You are connected",
				false: "Method not Allowed",
			}[IsCookieValid(w, r)],
		}}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	Resp = Response{}
	email := strings.TrimSpace(r.FormValue("email"))
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")
	if email != "" && username != "" && password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err == nil {
			_, err := db.CreateUser("", "user", username, email, "", string(password))
			if err == nil {
				SetSession(w, username)
			} else {
				Resp.Msg = append(Resp.Msg, err.Error())
				Resp.Action = "GetSignup();"
			}
		} else if err == bcrypt.ErrPasswordTooLong {
			Resp.Msg = append(Resp.Msg, "The password is too long!")
			Resp.Action = "GetSignup();"
		} else {
			Resp.Msg = append(Resp.Msg, err.Error())
			Resp.Action = "GetSignup();"
		}
	} else {
		Resp.Msg = append(Resp.Msg, "All fields are required!")
		Resp.Action = "GetSignup();"
	}

	if Resp.Msg == nil {
		Resp.Msg = append(Resp.Msg, "Account created successfully!")
	} else {
		Resp.Form.Username = username
		Resp.Form.Email = email
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
