#!/bin/bash

set -eu

# USAGE NOTES:
#
# - If the API URL is an IP address, the IP must be present in the certificate
#   as a common name (NOT alternative name) or it won't be validated correctly
#   due to a bug in Python 2.7.9+: https://bugs.python.org/issue23239
# - For a more basic but faster test, decrease SQLMAP_LEVEL and/or SQLMAP_RISK

CREDENTIAL_ROOT="$(mktemp -d)"
AUTH_FILE="${CREDENTIAL_ROOT}/auth_file.pem"
# Defaults level and risk to maximum to increase test coverage
SQLMAP_LEVEL="${SQLMAP_LEVEL:-5}"
SQLMAP_RISK="${SQLMAP_RISK:-3}"

BASEDIR="$(dirname ${0})"

setup_certs() {
    client_ca_cert_path="${CREDENTIAL_ROOT}/client_ca_cert.pem"
    client_ca_key_path="${CREDENTIAL_ROOT}/client_ca_key.pem"
    client_cert_path="${CREDENTIAL_ROOT}/client.pem"
    client_key_path="${CREDENTIAL_ROOT}/client_key.pem"

    echo "${CLIENT_CA_CERT}" > ${CREDENTIAL_ROOT}/client_ca_cert.pem
    echo "${CLIENT_CA_KEY}" > ${CREDENTIAL_ROOT}/client_ca_key.pem

    "${BASEDIR}/generate_certs.py" \
        -outputPath "${CREDENTIAL_ROOT}" \
        -caCert "${client_ca_cert_path}" \
        -caKey "${client_ca_key_path}"

    cat "${client_key_path}" "${client_cert_path}" > "${AUTH_FILE}"
}

run_sqlmap() {
    url="${1}"
    method="${2:-GET}"
    data="${3:-}"
    log_file="$(mktemp)"

    # choose level 3 verbosity as it shows injected payloads
    # https://github.com/sqlmapproject/sqlmap/wiki/Usage#output-verbosity
    sqlmap_command="sqlmap
        -u ${url}
        -v 3
        --auth-file ${AUTH_FILE}
        --dbms ${DATABASE_TYPE}
        -H \"content-type: application/json\"
        --level ${SQLMAP_LEVEL}
        --risk ${SQLMAP_RISK}
        --batch
        --fresh-queries
        --flush-session
        --method ${method}"

    # Can't use --data with GET or it will convert the request to a POST
    if [[ ! "${method}" == "GET" ]]; then
        sqlmap_command+=" --data ${data}"
    fi

    ${sqlmap_command} | tee "${log_file}"

    if grep -q "retrieved:" "${log_file}"; then
        echo "Found an injection - exiting"
        exit 1
    fi

    if grep -q "all tested parameters appear to be not injectable" "${log_file}"; then
        echo "No injections found!"
    fi
}

run_tests() {
    password_cred="/sqlmap/password"
    ssh_cred="/sqlmap/ssh"

    # Must test PUT before GET to ensure the credentials exist for GET tests.
    # Doesn't test DELETE because that causes 404s and causes sqlmap problems.

    # PUT tests
    run_sqlmap "${API_URL}/api/v1/data" PUT \
        "{\"name\":\"${password_cred}\",\"type\":\"password\",\"value\":\"test-password-value\",\"overwrite\":false}"
    run_sqlmap "${API_URL}/api/v1/data" PUT \
        "{\"name\":\"${ssh_cred}\",\"type\":\"ssh\",\"value\":{\"public_key\":\"test-public-key\"},\"overwrite\":false}"

    # GET tests
    run_sqlmap "${API_URL}/api/v1/data?name=${password_cred}"
    run_sqlmap "${API_URL}/api/v1/data?name=${ssh_cred}"

    # POST tests
    run_sqlmap "${API_URL}/api/v1/data" POST \
        "{\"name\":\"${password_cred}\",\"type\":\"password\",\"parameters\":{\"length\":10},\"overwrite\":false}"
    run_sqlmap "${API_URL}/api/v1/data" POST \
        "{\"name\":\"${ssh_cred}\",\"type\":\"ssh\",\"overwrite\":false}"

    # FIND tests
    run_sqlmap "${API_URL}/api/v1/data?path=/sqlmap" # find by path
    run_sqlmap "${API_URL}/api/v1/data?name-like=sqlmap" # find by name like
#    run_sqlmap "${API_URL}/api/v1/data?paths=true" # find paths - disabled because it can take too long

    # Clean up after ourselves
    curl "${API_URL}/api/v1/data?name=${password_cred}" \
        --silent \
        -X DELETE \
        --cert "${AUTH_FILE}"
    curl "${API_URL}/api/v1/data?name=${ssh_cred}" \
        --silent \
        -X DELETE \
        --cert "${AUTH_FILE}"
}

main() {
    setup_certs
    run_tests
}

main
