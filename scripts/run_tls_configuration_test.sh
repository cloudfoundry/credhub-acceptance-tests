#!/bin/sh

set -eu

MINIMUM_VERSION="7"

check_nmap_version() {
  if ! command -v nmap >/dev/null 2>&1; then
    echo "Please install nmap ${MINIMUM_VERSION}"
    exit 1
  fi

  regex="s/.*\([0-9]\{1,\}\)\.[0-9]\{1,\}.*/\1/g"
  major_version="$(nmap --version | grep -i "nmap version" | sed "${regex}")"
  echo "major version: ${major_version}"

  if [ "${major_version}" -ne "${MINIMUM_VERSION}" ]; then
    echo "Please install nmap ${MINIMUM_VERSION}"
    exit 1
  fi
}

run_nmap() {
  host_regex="s/.*https\{0,1\}:\/\/\(.*\):[0-9]\{2,\}/\1/g"
  port_regex="s/.*https\{0,1\}:\/\/.*:\([0-9]\{2,\}\)/\1/g"

  host="$(echo "${API_URL}" | sed "${host_regex}")"
  port="$(echo "${API_URL}" | sed "${port_regex}")"
  nmap --script ssl-enum-ciphers -p "${port}" "${host}"
}

check_cipher_preference() {
  expected="${1}"
  output="${2}"

  cipher_preference_regex="s/.*cipher preference: \([a-z]\{1,\}\)/\1/g"
  cipher_preference="$(echo "${output}" | grep "cipher preference" | sed "${cipher_preference_regex}")"

  if ! echo "${output}" | grep -q "cipher preference: ${expected}"; then
    echo "Unexpected cipher preference: expected \"${expected}\" but was \"${cipher_preference}\""
    echo "Nmap output: $output"
    exit 1
  else
    echo "Found correct cipher preference"
  fi
}

check_cipher() {
 cipher="$1"
 output="$2"
 if ! echo "$output" | grep -q "$cipher"; then
    echo "Did not find cipher in list of supported ciphers: $cipher"
    echo "Nmap output: $output"
    exit 1
 fi
 echo "Found correct cipher $cipher"
}
check_ciphers() {
  check_cipher "TLS_DHE_RSA_WITH_AES_128_GCM_SHA256" "$output"
  check_cipher "TLS_DHE_RSA_WITH_AES_256_GCM_SHA384" "$output"
  check_cipher "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256" "$output"
  check_cipher "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384" "$output"

}
main() {
  echo "Checking cipher preference for ${API_URL}"
  output=$(run_nmap)
  check_cipher_preference server "$output"
  check_ciphers "$output"
}

main
