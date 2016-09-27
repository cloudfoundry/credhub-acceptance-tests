

Acceptance test for Credential Manager

### Run Tests locally

Create `config.json` file in the project directory with desired target API URL like so:
```sh
cat > config.json <<EOF
{
  "api_url": "https://TARGET_API_IP:TARGET_API_PORT"
}
EOF
```

Runs local CredHub testing via:

```sh
ginkgo -r integration/
```
