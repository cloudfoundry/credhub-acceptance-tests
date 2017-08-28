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

check_cipher_preference() {
  expected="${1}"

  host_regex="s/.*https\{0,1\}:\/\/\(.*\):[0-9]\{2,\}/\1/g"
  port_regex="s/.*https\{0,1\}:\/\/.*:\([0-9]\{2,\}\)/\1/g"

  host="$(echo "${API_URL}" | sed "${host_regex}")"
  port="$(echo "${API_URL}" | sed "${port_regex}")"

  output="$(nmap --script ssl-enum-ciphers -p "${port}" "${host}")"
  cipher_preference_regex="s/.*cipher preference: \([a-z]\{1,\}\)/\1/g"
  cipher_preference="$(echo "${output}" | grep "cipher preference" | sed "${cipher_preference_regex}")"

  if ! echo "${output}" | grep -q "cipher preference: ${expected}"; then
    echo "Unexpected cipher preference: expected \"${expected}\" but was \"${cipher_preference}\""
    exit 1
  else
    echo "Found correct cipher preference"
  fi
}

main() {
  echo "Checking cipher preference for ${API_URL}"
  check_cipher_preference "server"
}

main
