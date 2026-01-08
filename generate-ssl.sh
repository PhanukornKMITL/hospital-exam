#!/bin/bash

# Generate self-signed SSL certificate for local testing
mkdir -p ssl

openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ssl/private.key \
  -out ssl/certificate.crt \
  -subj "/C=TH/ST=Bangkok/L=Bangkok/O=Hospital/CN=hospital-a.api.co.th"

echo "✅ SSL certificate created successfully!"
echo "📁 Files created:"
echo "   - ssl/certificate.crt"
echo "   - ssl/private.key"
