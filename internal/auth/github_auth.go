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

func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := utils.GenerateRandomState()

	params := url.Values{}
	params.Add("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	params.Add("redirect_uri", "https://forum.shasor.fr/auth/github/callback")
	params.Add("scope", "user,user:email")
	params.Add("state", state)

	http.SetCookie(w, &http.Cookie{Name: "oauth_state", Value: state, HttpOnly: true, Secure: true})
	http.Redirect(w, r, GithubAuthURL+"?"+params.Encode(), http.StatusFound)
}

func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	handlers.Resp = handlers.Response{}
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

func GithubExchangeCodeForToken(code string) (string, error) {
	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	values.Add("client_secret", os.Getenv("GITHUB_CLIENT_SECRET"))
	values.Add("redirect_uri", "https://forum.shasor.fr/auth/github/callback")

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
	// Requête pour obtenir les emails de l'utilisateur
	reqEmails, _ := http.NewRequest("GET", GithubEmailsURL, nil)
	reqEmails.Header.Add("Authorization", "Bearer "+token)

	respEmails, err := http.DefaultClient.Do(reqEmails)
	if err != nil {
		return nil, err
	}
	defer respEmails.Body.Close()

	var emails []map[string]interface{}
	if err := json.NewDecoder(respEmails.Body).Decode(&emails); err != nil {
		return nil, err
	}

	var primaryEmail string
	for _, email := range emails {
		if primary, ok := email["primary"].(bool); ok && primary {
			primaryEmail = email["email"].(string)
			break
		}
	}

	// Requête pour obtenir les autres informations de l'utilisateur
	req, _ := http.NewRequest("GET", GithubUserInfoURL, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	// Ajouter l'email primaire aux informations de l'utilisateur
	userInfo["email"] = primaryEmail
	return userInfo, nil
}

func GithubCreateOrUpdateUser(userInfo map[string]interface{}) (*db.User, error) {
	var errMsg []string
	if userInfo["email"] == nil {
		errMsg = append(errMsg, "email")
	}
	if userInfo["login"] == nil {
		errMsg = append(errMsg, "login")
	}
	if userInfo["avatar_url"] == nil {
		errMsg = append(errMsg, "avatar_url")
	}
	if errMsg != nil {
		str1 := "Unable to retrieve certain information :"
		str2 := "You have not been logged in."
		return nil, fmt.Errorf("%v %v. %v", str1, strings.Join(errMsg, ", "), str2)
	}
	var user *db.User
	if !db.UserExistByEmail(userInfo["email"].(string)) || !db.UserExistByUsername(userInfo["login"].(string)) {
		picture, err := utils.GetFileFromURL(userInfo["avatar_url"].(string))
		if err != nil {
			panic(err)
		}
		user, err = db.CreateUser("github", "user", userInfo["login"].(string), userInfo["email"].(string), picture, "")
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}
	} else {
		if db.UserExistByEmail(userInfo["email"].(string)) {
			user, _ = db.SelectUserByEmail(userInfo["email"].(string))
		} else {
			user, _ = db.SelectUserByUsername(userInfo["login"].(string))
		}
		if user.Provider == "github" {
			return user, nil
		}
		str := "The email address or username already exists!"
		return nil, fmt.Errorf("%v", str)
	}
	return user, nil
}
