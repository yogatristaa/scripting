#!/usr/bin/env python3
import jwt
import time
import sys

# Get PEM file path
if len(sys.argv) > 1:
    pem = sys.argv[1]
else:
    print(f"GitHub APP PEM File Not Specified")

# Get the App ID
if len(sys.argv) > 2:
    app_id = sys.argv[2]
else:
    (f"GitHub APP ID Not Specified")

# Open PEM
with open(pem, 'rb') as pem_file:
    signing_key = pem_file.read()

payload = {
    # Issued at time
    'iat': int(time.time()),
    # JWT expiration time (10 minutes maximum)
    'exp': int(time.time()) + 300,
    # GitHub App's identifier
    'iss': app_id
}

# Create JWT
encoded_jwt = jwt.encode(payload, signing_key, algorithm='RS256')

print(f"{encoded_jwt}")