package auth

import (
	"encoding/json"
	"forum/internal/db"
	"forum/internal/handlers"
	"forum/internal/utils"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := utils.GenerateRandomState()

	params := url.Values{}
	params.Add("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	params.Add("redirect_uri", "https://localhost:8080/auth/github/callback")
	params.Add("scope", "user:email")
	params.Add("state", state)

	http.SetCookie(w, &http.Cookie{Name: "oauth_state", Value: state, HttpOnly: true, Secure: true})
	http.Redirect(w, r, GithubAuthURL+"?"+params.Encode(), http.StatusFound)
}

func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")

	oauthCookie, err := r.Cookie("oauth_state")
	if err != nil || state != oauthCookie.Value {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	token, err := GithubExchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	userInfo, err := GithubGetUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	user, err := GithubCreateOrUpdateUser(userInfo)
	if err != nil {
		http.Error(w, "Failed to create or update user", http.StatusInternalServerError)
		return
	}

	// Créer une session pour l'utilisateur
	handlers.ClearSession(w, r, "oauth_state")
	handlers.SetSession(w, user.Username)
	// ...

	http.Redirect(w, r, "/", http.StatusFound)
}

func GithubExchangeCodeForToken(code string) (string, error) {
	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	values.Add("client_secret", os.Getenv("GITHUB_CLIENT_SECRET"))
	values.Add("redirect_uri", "https://localhost:8080/auth/github/callback")

	req, err := http.NewRequest("POST", GithubTokenURL, strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")

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

func GithubGetUserInfo(token string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", GithubUserInfoURL, nil)
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

func GithubCreateOrUpdateUser(userInfo map[string]interface{}) (*db.User, error) {
	var user *db.User
	if !db.UserExistByEmail(userInfo["email"].(string)) {
		picture, err := utils.GetFileFromURL(userInfo["avatar_url"].(string))
		if err != nil {
			panic(err)
		}
		user, err = db.CreateUser("user", userInfo["login"].(string), userInfo["email"].(string), picture, "")
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		user, err = db.SelectUserByEmail(userInfo["email"].(string))
		if err != nil {
			panic(err)
		}
	}
	return user, nil
}
