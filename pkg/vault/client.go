package vault

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

type Config struct {
	Environment string
	ServerPort  string
	LogLevel    string
	EnableCORS  bool
	JWTSecret   string
	APIKey      string
}

type Client struct {
	client *api.Client
	env    string
}

func NewClient(environment string) (*Client, error) {
	config := api.DefaultConfig()
	
	// Get Vault address from environment
	config.Address = os.Getenv("VAULT_ADDR")
	if config.Address == "" {
		return nil, fmt.Errorf("VAULT_ADDR environment variable is required")
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %w", err)
	}

	// Get token from environment
	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("VAULT_TOKEN environment variable is required")
	}
	client.SetToken(token)

	return &Client{
		client: client,
		env:    environment,
	}, nil
}

// LoadConfig loads configuration from Vault
func (v *Client) LoadConfig() (*Config, error) {
	cfg := &Config{}

	// Load app configuration
	appSecret, err := v.client.KVv2("secret").Get(
		context.Background(),
		fmt.Sprintf("learngo/%s/config", v.env),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read app config from Vault: %w", err)
	}

	if appSecret != nil {
		appData := appSecret.Data
		cfg.Environment = getString(appData, "environment")
		cfg.ServerPort = getString(appData, "server_port")
		cfg.LogLevel = getString(appData, "log_level")
		cfg.EnableCORS = getBool(appData, "enable_cors")
	}

	// Load auth configuration
	authSecret, err := v.client.KVv2("secret").Get(
		context.Background(),
		fmt.Sprintf("learngo/%s/auth", v.env),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read auth config from Vault: %w", err)
	}

	if authSecret != nil {
		authData := authSecret.Data
		cfg.JWTSecret = getString(authData, "jwt_secret")
		cfg.APIKey = getString(authData, "api_key")
	}

	return cfg, nil
}

// Helper functions
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok && val != nil {
		return val.(string)
	}
	return ""
}

func getBool(data map[string]interface{}, key string) bool {
	if val, ok := data[key]; ok && val != nil {
		return val.(bool)
	}
	return false
}
