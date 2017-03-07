# Acceptance test suite for CredHub

CredHub manages credentials like passwords, certificates, ssh keys, rsa keys, strings (arbitrary values) and CAs. CredHub provides a CLI and API to get, set, generate and securely store such credentials.

* [CredHub Tracker][1]
[1]:https://www.pivotaltracker.com/n/projects/1977341

See additional repos for more info:

* [credhub](https://github.com/cloudfoundry-incubator/credhub) :     CredHub server code 
* [credhub-cli](https://github.com/cloudfoundry-incubator/credhub-cli) :     command line interface for credhub
* [credhub-release](https://github.com/pivotal-cf/credhub-release) : BOSH release of CredHub server

### Get prerequisites

Ensure that you have a local version of the CredHub CLI checked out in your $GOPATH

```sh
go get github.com/cloudfoundry-incubator/credhub-cli
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
