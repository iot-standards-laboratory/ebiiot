#!/bin/bash
set -e

read -p 'Enter domain (default: localhost): ' domainName
domainName=${domainName:-localhost}

echo "Generating CA key and certificate:"
openssl req -x509 -sha256 -nodes -days 3650 -newkey rsa:2048 \
  -keyout ca.key -out ca.pem \
  -subj "/O=Godopu Certificate Authority/"

echo "Generating CSR"
openssl req -out cert.csr -new -newkey rsa:2048 -nodes -keyout priv.key \
  -subj "/O=Godopu/"

echo "Sign certificate:"
openssl x509 -req -sha256 -days 3650 -in cert.csr  -out cert.pem \
  -CA ca.pem -CAkey ca.key -CAcreateserial \
  -extfile <(printf "subjectAltName=DNS:$domainName")

# debug output the certificate
openssl x509 -noout -text -in cert.pem

# make .crt file from .pem
# openssl x509 -in -ca.pem -inform pem -out ca.crt

# we don't need the CA key, the serial number and the CSR any more
rm ca.key cert.csr

