package auth

import (
	"encoding/json"
	"fmt"
	"forum/internal/db"
	"forum/internal/handlers"
	"forum/internal/utils"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func DiscordLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := utils.GenerateRandomState()

	params := url.Values{}
	params.Add("client_id", os.Getenv("DISCORD_CLIENT_ID"))
	params.Add("redirect_uri", "https://localhost:8080/auth/discord/callback")
	params.Add("scope", "email identify")
	params.Add("response_type", "code")
	params.Add("state", state)

	http.SetCookie(w, &http.Cookie{Name: "oauth_state", Value: state, HttpOnly: true, Secure: true})
	http.Redirect(w, r, DiscordAuthURL+"?"+params.Encode(), http.StatusFound)
}

func DiscordCallbackHandler(w http.ResponseWriter, r *http.Request) {
	handlers.Resp = handlers.Response{}
	code := r.FormValue("code")
	state := r.FormValue("state")

	oauthCookie, err := r.Cookie("oauth_state")
	if err != nil || state != oauthCookie.Value {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}
	token, err := DiscordExchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	userInfo, err := DiscordGetUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	user, err := DiscordCreateOrUpdateUser(userInfo)
	if err != nil {
		handlers.Resp.Msg = append(handlers.Resp.Msg, err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Créer une session pour l'utilisateur
	handlers.ClearSession(w, r, "oauth_state")
	handlers.SetSession(w, user.Username)
	// ...

	http.Redirect(w, r, "/", http.StatusFound)
}

func DiscordExchangeCodeForToken(code string) (string, error) {
	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", os.Getenv("DISCORD_CLIENT_ID"))
	values.Add("client_secret", os.Getenv("DISCORD_CLIENT_SECRET"))
	values.Add("redirect_uri", "https://localhost:8080/auth/discord/callback")
	values.Add("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", DiscordTokenURL, strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

func DiscordGetUserInfo(token string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", DiscordUserInfoURL, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}

func DiscordCreateOrUpdateUser(userInfo map[string]interface{}) (*db.User, error) {
	var errMsg []string
	if userInfo["email"] == nil {
		errMsg = append(errMsg, "email")
	}
	if userInfo["global_name"] == nil {
		errMsg = append(errMsg, "global_name")
	}
	if userInfo["id"] == nil {
		errMsg = append(errMsg, "id")
	}
	if userInfo["avatar"] == nil {
		errMsg = append(errMsg, "avatar")
	}
	if errMsg != nil {
		str1 := "Unable to retrieve certain information :"
		str2 := "You have not been logged in."
		return nil, fmt.Errorf("%v %v. %v", str1, strings.Join(errMsg, ", "), str2)
	}
	var user *db.User
	if !db.UserExistByEmail(userInfo["email"].(string)) || !db.UserExistByUsername(userInfo["global_name"].(string)) {
		picture, err := utils.GetFileFromURL(DiscordAvatarURL + "/" + userInfo["id"].(string) + "/" + userInfo["avatar"].(string))
		if err != nil {
			panic(err)
		}
		user, err = db.CreateUser("discord", "user", userInfo["global_name"].(string), userInfo["email"].(string), picture, "")
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}
	} else {
		if db.UserExistByEmail(userInfo["email"].(string)) {
			user, _ = db.SelectUserByEmail(userInfo["email"].(string))
		} else {
			user, _ = db.SelectUserByUsername(userInfo["global_name"].(string))
		}
		if user.Provider == "discord" {
			return user, nil
		}
		str := "The email address or username already exists!"
		return nil, fmt.Errorf("%v", str)
	}
	return user, nil
}
