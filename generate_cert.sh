#!/bin/bash

# Directory where the certificate files should be stored
CERT_DIR="./certs"
CERT_FILE="$CERT_DIR/server.crt"
KEY_FILE="$CERT_DIR/server.key"

# Check if the certificate directory exists, create it if not
if [ ! -d "$CERT_DIR" ]; then
    mkdir -p "$CERT_DIR"
fi

# Check if the certificate and key already exist
if [ ! -f "$CERT_FILE" ] || [ ! -f "$KEY_FILE" ]; then
    echo "Certificate or key not found. Generating new certificate..."
    # Generate the certificate and key using OpenSSL
    if openssl req -x509 -nodes -newkey rsa:2048 -keyout "$KEY_FILE" -out "$CERT_FILE" -config openssl.conf -extensions v3_req -days 365 &> /dev/null; then
        echo "Certificate and key generated at $CERT_FILE and $KEY_FILE"
    else
        echo "SSL certificate generation error"
    fi
else
    echo "Certificate and key already exist. Skipping generation"
fi