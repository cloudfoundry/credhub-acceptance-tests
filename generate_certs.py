#!/usr/bin/env python

import argparse
import subprocess
import os
import sys
import uuid

cert_dir = os.path.join(os.getcwd(), 'certs/')

def setup_cert_dir():
    print "initializing certificate directory"
    os.system("rm -rf certs")
    os.makedirs("certs")

def generate_valid_cert(ca_cert_path, ca_key_path):
    generate_cert(ca_cert_path, ca_key_path, "client", "30")

def generate_self_signed_cert(file_base_name, days):
    cert_path = os.path.join(cert_dir, file_base_name+".pem")
    key_path = os.path.join(cert_dir, file_base_name+"_key.pem")
    subprocess.call(["openssl", "req", "-x509", "-newkey", "rsa:2048", "-days", days, "-subj", "/CN=credhub_test_client/OU=app:20045e7a-4686-48ba-b523-b5378357073f",
                     "-nodes", "-sha256", "-keyout", key_path, "-out", cert_path])

def generate_cert(ca_cert_path, ca_key_path, file_base_name, days):
    cert_path = os.path.join(cert_dir, file_base_name+".pem")
    key_path = os.path.join(cert_dir, file_base_name+"_key.pem")
    client_csr_path = os.path.join(cert_dir, file_base_name+".csr")
    # generate keypair
    subprocess.call(["openssl", "genrsa", "-out", key_path, "2048"])

    # create CSR
    app_uuid = uuid.uuid4()
    subprocess.call(["openssl", "req", "-new", "-key", key_path, "-out", client_csr_path, "-subj", "/CN=credhub_test_client/OU=app:" + str(app_uuid)])

    # generate client certificate
    subprocess.call(["openssl", "x509", "-req", "-in", client_csr_path, "-extfile", "openssl_extra.cnf", "-extensions", "v3_req", "-CA", ca_cert_path, "-CAkey", ca_key_path,
                     "-CAcreateserial", "-days", days, "-sha256", "-out", cert_path])

def make_unknown_cert_and_ca():
    # Make a new CA
    unknown_ca_path = os.path.join(cert_dir, "unknown_ca.pem")
    unknown_ca_key_path = os.path.join(cert_dir, "unknown_ca_key.pem")
    subprocess.call(["openssl", "req", "-x509", "-newkey", "rsa:2048", "-days", "365", "-sha256", "-nodes", "-subj", "/CN=credhub_client_ca/OU=app:20045e7a-4686-48ba-b523-b5378357073f", "-keyout", unknown_ca_key_path, "-out", unknown_ca_path])

    # Sign a new cert with unknown CA
    generate_cert(unknown_ca_path, unknown_ca_key_path, "unknown", "30")

def generate_bad_certs(ca_cert_path, ca_key_path):
    generate_self_signed_cert("selfsigned", "30")
    make_unknown_cert_and_ca()
    generate_cert(ca_cert_path, ca_key_path, "expired", "-1")

tool_desc = "TLS certificate generator for CredHub acceptance tests"
parser = argparse.ArgumentParser(description=tool_desc)

parser.add_argument('-caCert','-c','--c', dest='ca_cert_path', help="Path to PEM encoded CA public cert")
parser.add_argument('-caKey','-k','--k', dest='ca_key_path', help="Path to PEM encoded CA private key")

args = parser.parse_args()

if not (args.ca_cert_path and args.ca_key_path):
    parser.print_usage()
    sys.exit(1)

setup_cert_dir()

generate_valid_cert(args.ca_cert_path, args.ca_key_path)

generate_bad_certs(args.ca_cert_path, args.ca_key_path)
