# Development app config
VALT_TOKEN="root-token-123"
curl \
  --header "X-Vault-Token: $VALT_TOKEN" \
  --request POST \
  --data '{
    "data": {
      "docker-password": "sq@1234",
      "VAULT_ADDR": "http://vault.qwerfvcxza.site",
      "DOCKER_REGISTRY": "chouleang"
    }
  }' \
  http://vault.qwerfvcxza.site/v1/secret/data/learngo/jenkins/go-operator

