package auth

const (
	GoogleAuthURL     = "https://accounts.google.com/o/oauth2/auth"
	GoogleTokenURL    = "https://oauth2.googleapis.com/token"
	GoogleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"

	GithubAuthURL     = "https://github.com/login/oauth/authorize"
	GithubTokenURL    = "https://github.com/login/oauth/access_token"
	GithubUserInfoURL = "https://api.github.com/user"
	GithubEmailsURL   = "https://api.github.com/user/emails"

	DiscordAuthURL     = "https://discord.com/oauth2/authorize"
	DiscordTokenURL    = "https://discord.com/api/oauth2/token"
	DiscordUserInfoURL = "https://discord.com/api/users/@me"
	DiscordAvatarURL   = "https://cdn.discordapp.com/avatars"
)
