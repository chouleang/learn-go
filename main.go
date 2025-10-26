package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go-hello-operator/pkg/vault" // Import your vault package
)

func main() {
	// Determine environment (development, production, etc.)
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // Default to development
	}

	// Initialize Vault client
	vaultClient, err := vault.NewClient(env)
	if err != nil {
		log.Fatalf("❌ Failed to initialize Vault client: %v", err)
	}

	// Load configuration from Vault
	cfg, err := vaultClient.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration from Vault: %v", err)
	}

	log.Printf("✅ Loaded configuration for: %s environment", cfg.Environment)
	log.Printf("📍 Server port: %s", cfg.ServerPort)
	log.Printf("📝 Log level: %s", cfg.LogLevel)
	log.Printf("🌐 CORS enabled: %v", cfg.EnableCORS)
	log.Printf("🔐 JWT secret loaded: %t", cfg.JWTSecret != "")
	log.Printf("🔑 API key loaded: %t", cfg.APIKey != "")

	// Setup HTTP routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Learn Go! 🎯\n")
		fmt.Fprintf(w, "Environment: %s\n", cfg.Environment)
		fmt.Fprintf(w, "All configuration loaded securely from Vault! ✅")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "healthy", "environment": "%s"}`, cfg.Environment)
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"environment": "%s",
			"server_port": "%s",
			"log_level": "%s",
			"cors_enabled": %v,
			"jwt_configured": %v,
			"api_key_configured": %v
		}`, cfg.Environment, cfg.ServerPort, cfg.LogLevel,
			cfg.EnableCORS, cfg.JWTSecret != "", cfg.APIKey != "")
	})

	// Protected endpoint that uses the API key from Vault
	http.HandleFunc("/api/secure", func(w http.ResponseWriter, r *http.Request) {
		// Check API key from header
		providedKey := r.Header.Get("X-API-Key")
		if providedKey != cfg.APIKey {
			http.Error(w, `{"error": "Unauthorized - Invalid API Key"}`, http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"message": "Access granted!",
			"environment": "%s",
			"secret_data": "This is protected information"
		}`, cfg.Environment)
	})

	// Start server with port from Vault configuration
	log.Printf("🚀 Starting server on :%s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, nil))
}
