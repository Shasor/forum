package server

import (
	"bufio"
	"fmt"
	"forum/internal/db"
	"forum/internal/router"
	"forum/internal/security"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func InitServer() {
	// Initialize the database
	DB := db.GetDB()
	defer DB.Close()

	// Set up routes
	mux := router.SetupRoutes(DB)

	// Load custom TLS configuration
	tlsConfig := security.LoadTLSConfig()

	// Create the HTTP server
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		MaxHeaderBytes:    1 << 20,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       30 * time.Second,
		WriteTimeout:      10 * time.Second,
		TLSConfig:         tlsConfig,
	}

	// Start the server with HTTPS
	fmt.Println("Server started at: https://localhost:8080")
	log.Fatal(server.ListenAndServeTLS("certs/server.crt", "certs/server.key"))
}

func NewServer(port string, readTimeout, writeTimeout, idleTimeout, readHeaderTimeout time.Duration, maxHeaderBytes int) *http.Server {
	return &http.Server{
		Addr:              port,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
}

func loadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	return scanner.Err()
}
