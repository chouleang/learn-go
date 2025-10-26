# production app config
VALT_TOKEN="root-token-123"
curl \
  --header "X-Vault-Token: $VALT_TOKEN" \
  --request POST \
  --data '{
    "data": {
      "environment": "production",
      "server_port": "8080",
      "log_level": "warn",
      "enable_cors": false,
      "api_timeout": "10s",
      "cache_ttl": "1m"
    }
  }' \
  http://do.dns.test:8200/v1/secret/data/learngo/production/config

# production auth config
curl \
  --header "X-Vault-Token: $VALT_TOKEN" \
  --request POST \
  --data '{
    "data": {
      "jwt_secret": "prod-jwt-secret-1234567890",
      "session_key": "prod-session-key-123456",
      "token_expiry": "1h"
    }
  }' \
  http://do.dns.test:8200/v1/secret/data/learngo/production/auth
