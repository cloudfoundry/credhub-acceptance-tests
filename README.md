# Acceptance test for CredHub

# Get prerequisites

Ensure that you have a local version of the CredHub CLI checked out in your $GOPATH

```sh
go get github.com/pivotal-cf/credhub-cli
```

### Run Tests locally

Target your local API by running:

```sh
cat <<EOF > config.json
{
  "api_url": "https://${YOUR_IP_HERE}:8844",
  "api_username":"${YOUR_USERNAME}",
  "api_password":"${YOUR_PASSWORD}"
}
EOF
```

Runs local CredHub testing via:

```sh
./run_tests.sh
```
