package handlers

import (
	"fmt"
	"forum/internal/db"
	"net/http"
	"strconv"
)

func RoleHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST and if the user is connected (cookie is valid)
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
	role := r.FormValue("role")
	otherID, err := strconv.Atoi(r.FormValue("otherID"))
	if err != nil {
		panic(err)
	}

	if role != "user" && role != "moderator" {
		panic("The role must be equal to 'user' or 'moderator'")
	}

	_, err = db.SelectUserByID(otherID)
	if err == nil {
		err = db.UpdateUserRole(otherID, role)
		if err != nil {
			Resp.Msg = append(Resp.Msg, fmt.Sprintf("%v", err))
		}
	} else {
		Resp.Msg = append(Resp.Msg, fmt.Sprintf("%v", err))
	}

	if Resp.Msg == nil {
		Resp.Msg = append(Resp.Msg, "The User Role has been successfully updated!")
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
