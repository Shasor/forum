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
	"os/exec"
	"strings"
	"time"
)

func InitServer() {
	if err := loadEnv(".env"); err != nil {
		fmt.Printf("Error load .env: %v\n", err)
	}
	// Initialize the database
	DB := db.GetDB()
	defer DB.Close()

	// Set up routes
	mux := router.SetupRoutes(DB)

	// Check and create the certificate if necessary
	err := checkAndCreateCert("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatalf("Error generating certificate: %v", err)
	}

	// Load custom TLS configuration
	tlsConfig := security.LoadTLSConfig()

	// Create the HTTP server
	server := NewServer(
		":8080",
		10*time.Second, //Read timeout
		10*time.Second, //Write timeout
		30*time.Second, //Idle timeout
		10*time.Second, //Read header timeout
		1<<20,          //Max header bytes
	)
	server.Handler = mux
	server.TLSConfig = tlsConfig

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

func checkAndCreateCert(certFile, keyFile string) error {
	// Check if the certificate file exists
	_, err := os.Stat(certFile)
	if os.IsNotExist(err) {
		// If the certificate does not exist, create it using the shell script
		fmt.Println("Certificate not found, generating new one...")
		return generateCertFromScript()
	}

	// Check if the key file exists
	_, err = os.Stat(keyFile)
	if os.IsNotExist(err) {
		// If the key file does not exist, create it using the shell script
		fmt.Println("Key file not found, generating new one...")
		return generateCertFromScript()
	}

	// If both the certificate and key exist, return nil (no error)
	return nil
}

func generateCertFromScript() error {
	// Execute the generate_cert.sh script
	cmd := exec.Command("./generate_cert.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the script and wait for it to finish
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute generate_cert.sh: %v", err)
	}

	return nil
}
