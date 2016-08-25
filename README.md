Acceptance test for Credential Manager

### Run Tests locally

Create `config/config.json` file with desired target API URL like so:
```sh
cat > config/config.json <<EOF
{
  "api_url": "https://TARGET_API_IP:TARGET_API_PORT"
}
EOF
```

Runs local CredHub
