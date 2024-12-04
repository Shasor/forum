package auth

const (
	GoogleAuthURL     = "https://accounts.google.com/o/oauth2/auth"
	GoogleTokenURL    = "https://oauth2.googleapis.com/token"
	GoogleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"

	GithubAuthURL     = "https://github.com/login/oauth/authorize"
	GithubTokenURL    = "https://github.com/login/oauth/access_token"
	GithubUserInfoURL = "https://api.github.com/user"

	DiscordAuthURL = "https://discord.com/oauth2/authorize?client_id=1313873177348280420&response_type=code&redirect_uri=https%3A%2F%2Flocalhost%3A8080%2Fauth%2Fdiscord%2Fcallback&scope=email+identify"
	DiscordTokenURL = "https://discord.com/api/oauth2/token"
	DiscordUserInfoURL = "https://discord.com/api/users/@me"
)
