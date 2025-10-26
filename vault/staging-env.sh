# staging app config
VALT_TOKEN="root-token-123"
curl \
  --header "X-Vault-Token: $VALT_TOKEN" \
  --request POST \
  --data '{
    "data": {
      "environment": "staging",
      "server_port": "8080",
      "log_level": "debug",
      "enable_cors": true,
      "api_timeout": "15s",
      "cache_ttl": "2m"
    }
  }' \
  http://do.dns.test:8200/v1/secret/data/learngo/staging/config

# staging auth config
curl \
  --header "X-Vault-Token: $VALT_TOKEN" \
  --request POST \
  --data '{
    "data": {
      "jwt_secret": "staging-jwt-secret-1234567890",
      "session_key": "staging-session-key-123456",
      "token_expiry": "12h"
    }
  }' \
  http://do.dns.test:8200/v1/secret/data/learngo/staging/auth
