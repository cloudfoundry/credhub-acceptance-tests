# Acceptance test suite for CredHub
test
CredHub manages credentials like passwords, certificates, ssh keys, rsa keys, strings (arbitrary values) and CAs. CredHub provides a CLI and API to get, set, generate and securely store such credentials.

* [CredHub Tracker](https://www.pivotaltracker.com/n/projects/1977341)

See additional repos for more info:

* [credhub](https://github.com/cloudfoundry-incubator/credhub) :     CredHub server code 
* [credhub-cli](https://github.com/cloudfoundry-incubator/credhub-cli) :     command line interface for credhub
* [credhub-release](https://github.com/pivotal-cf/credhub-release) : BOSH release of CredHub server

### Get prerequisites

Ensure that you have a local version of the CredHub CLI and ginkgo checked out in your $GOPATH

Install the CredHub CLI
```sh
go get code.cloudfoundry.org/credhub-cli
```

To install ginkgo see [ginkgo installation](https://github.com/onsi/ginkgo#global-installation)


### Run Tests locally

Target your local API by running:

```sh
cat <<EOF > test_config.json
{
  "api_url": "https://${YOUR_IP_HERE}:8844",
  "api_username":"${YOUR_USERNAME}",
  "api_password":"${YOUR_PASSWORD}",
  "credential_root":"${YOUR_CREDHUB_CA_PATH}",
  "uaa_ca":"${UAA_CA_PEM_FILE}"
}
EOF
```

Runs local CredHub testing via:

```sh
./scripts/run_tests.sh
```

To run with a locally built credhub-cli you can replace the build step in 
[the before suite](https://github.com/cloudfoundry-incubator/credhub-acceptance-tests/blob/main/integration_test/integration_suite_test.go#L59)
with the path to your CLI.

### Run Application Smoke Tests

Target your desired environment:

```sh
cat <<EOF > test_config.json
{
  "api_url": "https://${YOUR_IP_HERE}:8844",
  "api_username":"${YOUR_USERNAME}",
  "api_password":"${YOUR_PASSWORD}"
}
EOF
```

Run smoke test suite via:

```sh
./scripts/run_smoke_tests.sh
```
