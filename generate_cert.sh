#!/bin/bash

# Directory where the certificate files should be stored
CERT_DIR="./certs"

# Paths to the certificate and private key
CERT_FILE="$CERT_DIR/server.crt"
KEY_FILE="$CERT_DIR/server.key"

# Check if the certificate directory exists, create it if not
if [ ! -d "$CERT_DIR" ]; then
    mkdir -p "$CERT_DIR"
    echo "Created certificate directory: $CERT_DIR"
fi

# Check if the certificate and key already exist
if [ ! -f "$CERT_FILE" ] || [ ! -f "$KEY_FILE" ]; then
    echo "Certificate or key not found. Generating new certificate..."

    # Create the OpenSSL configuration file dynamically if needed
    cat > openssl.conf <<EOL
[ req ]
distinguished_name = req_distinguished_name
x509_extensions = v3_req
prompt = no

[ req_distinguished_name ]
C = FR
ST = Île-de-France
L = Rouen
O = Forum
OU = Zone01
CN = localhost
emailAddress = contact@monsite.fr

[ v3_req ]
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = localhost
EOL

    # Generate the certificate and key using OpenSSL
    openssl req -x509 -nodes -newkey rsa:2048 -keyout "$KEY_FILE" -out "$CERT_FILE" -config openssl.conf -extensions v3_req -days 365
    echo "Certificate and key generated at $CERT_FILE and $KEY_FILE"
else
    echo "Certificate and key already exist. Skipping generation."
fi