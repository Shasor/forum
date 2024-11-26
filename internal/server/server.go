package server

import (
	"fmt"
	"forum/internal/handlers"
	"forum/internal/middlewares"
	"log"
	"time"
)

func InitServer() {
	err := checkAndCreateCert("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatalf("Error generating certificate: %v", err)
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

	server.Use(middlewares.NotFoundMiddleware)
	server.Use(middlewares.RecoverMiddleware)
	server.Use(middlewares.NewRateLimiterMap().RateLimitMiddleware)

	if err := server.Start(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
