#!/bin/bash

set -ex

# add root ca
cat >> /etc/ssl/certs/ca-certificates.crt <<EOF
-----BEGIN CERTIFICATE-----
MIICxTCCAa2gAwIBAgIUDmd16kLGLKfEoGDdXwV6uvfU06kwDQYJKoZIhvcNAQEL
BQAwEjEQMA4GA1UEAwwHZGVmYXVsdDAeFw0xNjA4MjkyMjA4MzFaFw0xNzA4Mjky
MjA4MzFaMBIxEDAOBgNVBAMMB2RlZmF1bHQwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQCxzvwAVJSHnKkFjK86JIKzEZixUGQ97qb4VdYpZ30+i/NBmbne
/PJ0IXhhRfRamiu0wFSkDadfpTGufMW2xsH2Z2ob4t2RhXHnxPl9yceaonsL6mpg
f0tr4LU/MtodZk3cNWdUTexgdn9IL1uaFyPN2RPmew6wsadtOGSK+w2gxU0WGojp
8iBwOBkNlBvl28oKdyBI2/hQw/zKLq9VUFJCbapXYpCpJD7gs3NZ6NJzA5T14c11
p4dgEEcMS4+IwrMalGdi+2NXkjaV4aYb18xk3jKuvU8MIZJfdGpFf0H9hJy1aIIm
lDwFX/bueme44uKoAjVQRpCtYe2jSrdlgbIvAgMBAAGjEzARMA8GA1UdEwEB/wQF
MAMBAf8wDQYJKoZIhvcNAQELBQADggEBADtrg3mNjvAu3Z5+ivQFu3ETlxFDENoG
P9tUKped9d9J0vG5wbclRnFSlTT7uv6t6pOd6arz38DnXAI1OqGIqXr/EKNYmEs/
EU+zjx0Ku/bX3kwZYispRv29GS0k0E0H0F2WVVNWnBZx9sfXO/HUrmlw3tfSG9T8
rLntxsJMMn7C4SebHk0nwWEomA8P0aQjD2NAIy1H3ucGskwdgVqhbKCarL/6WtKx
Ng8ZGqVKZ0K9dFVqU/7LQPDDu6gCOe9c982mC4sRltKlQbmvwgzNZUuyRzIV4COu
1yov6e0CKqyhb4d1KTfECCkol42mP1CaelxJnedJ9+5yZTdIv1K4ylI=
-----END CERTIFICATE-----
EOF

export GOPATH=$PWD/go
export PATH=$PATH:$GOPATH/bin

cd go/src/github.com/pivotal-cf/credhub-cli
make dependencies
cd ../cred-hub-acceptance-tests
cat > config.json <<EOF
{
  "api_url": "$API_URL"
}
EOF

ginkgo -r integration
