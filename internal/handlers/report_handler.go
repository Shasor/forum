package handlers

import (
	"fmt"
	"forum/internal/db"
	"net/http"
	"strconv"
	"time"
)

func ReportHandler(w http.ResponseWriter, r *http.Request) {
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

	postID, err := strconv.Atoi(r.FormValue("postID"))
	if err != nil {
		panic(err)
	}

	user := GetUserFromCookie(w, r)

	post, err := db.SelectPostByID(postID)
	if err != nil {
		Resp.Msg = append(Resp.Msg, err.Error())
	}

	err = db.AddNotification("report", date, user.ID, 0, post.ID, post.ParentID)
	if err != nil {
		Resp.Msg = append(Resp.Msg, err.Error())
	}

	if Resp.Msg == nil {
		Resp.Msg = append(Resp.Msg, "Your message has been sent to the moderators!")
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
