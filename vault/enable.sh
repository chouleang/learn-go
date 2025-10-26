curl \
  --header "X-Vault-Token: root-token-123" \
  --request POST \
  --data '{"type":"kv", "options":{"version":"2"}}' \
  http://do.dns.test:8200/v1/sys/mounts/secret
