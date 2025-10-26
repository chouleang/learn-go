# Development app config
VALT_TOKEN="root-token-123"
curl \
  --header "X-Vault-Token: $VALT_TOKEN" \
  --request POST \
  --data '{
    "data": {
      "environment": "development",
      "server_port": "8080",
      "log_level": "debug",
      "enable_cors": true,
      "api_timeout": "30s",
      "cache_ttl": "5m"
    }
  }' \
  http://do.dns.test:8200/v1/secret/data/learngo/development/config

# Development auth config
curl \
  --header "X-Vault-Token: $VALT_TOKEN" \
  --request POST \
  --data '{
    "data": {
      "jwt_secret": "dev-jwt-secret-1234567890",
      "session_key": "dev-session-key-123456",
      "token_expiry": "24h"
    }
  }' \
  http://do.dns.test:8200/v1/secret/data/learngo/development/auth
