package server

import (
	"fmt"
	"forum/internal/auth"
	"forum/internal/db"
	"forum/internal/handlers"
	"forum/internal/middlewares"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func InitServer() {
	password, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.DefaultCost)
	if err == nil {
		_, _ = db.CreateUser("", "admin", "admin", "admin@admin", "", string(password))
	} else {
		fmt.Println("Error creating admin : ", err)
	}

	// Create the HTTP server
	server := NewServer(":8080", 10*time.Second, 10*time.Second, 30*time.Second, 10*time.Second, 1<<20)

	server.Handle("/", handlers.IndexHandler)
	server.Handle("/signup", handlers.SignupHandler)
	server.Handle("/login", handlers.LoginHandler)
	server.Handle("/create-post", handlers.CreatePostHandler)
	server.Handle("/react", handlers.ReactToPost)
	server.Handle("/follow", handlers.FollowHandler)
	server.Handle("/logout", handlers.LogoutHandler)
	server.Handle("/delete", handlers.DeleteHandler)
	server.Handle("/edit", handlers.EditProfileHandler)
	server.Handle("/users", handlers.GetUserInfo)
	server.Handle("/delete-post", handlers.DeletePostHandler)
	server.Handle("/role", handlers.RoleHandler)
	server.Handle("/modify-post", handlers.ModifyPostHandler)
	server.Handle("/report", handlers.ReportHandler)
	server.Handle("/request", handlers.RequestHandler)
	server.Handle("/notifications", handlers.NotificationHandler)
	server.Handle("/notifications/clear", handlers.NotificationClearHandler)
	server.Handle("/auth/google/login", auth.GoogleLoginHandler)
	server.Handle("/auth/google/callback", auth.GoogleCallbackHandler)
	server.Handle("/auth/github/login", auth.GithubLoginHandler)
	server.Handle("/auth/github/callback", auth.GithubCallbackHandler)
	server.Handle("/auth/discord/login", auth.DiscordLoginHandler)
	server.Handle("/auth/discord/callback", auth.DiscordCallbackHandler)

	server.Use(middlewares.NotFoundMiddleware)
	server.Use(middlewares.RecoverMiddleware)
	server.Use(middlewares.NewRateLimiterMap().RateLimitMiddleware)

	if err := server.Start(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
