package server

import (
	"fmt"
	"forum/internal/handlers"
	"forum/internal/middlewares"
	"forum/internal/db"
	"log"
	"time"
	"golang.org/x/crypto/bcrypt"
)

func InitServer() {
	err := checkAndCreateCert("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatalf("Error generating certificate: %v", err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err ==nil{
		_, _ := db.CreateUser("admin", "admin", "admin@admin", "", string(password))
	}else {
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

	server.Use(middlewares.NotFoundMiddleware)
	server.Use(middlewares.RecoverMiddleware)
	server.Use(middlewares.NewRateLimiterMap().RateLimitMiddleware)

	if err := server.Start(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
