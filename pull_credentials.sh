#!/bin/bash
# pull_credentials.sh - Securely retrieve identity tokens from Google Drive
FILE_ID=$1
OUTPUT=".env"

if [ -z "$FILE_ID" ]; then
    echo "Usage: bash pull_credentials.sh <GOOGLE_DRIVE_FILE_ID>"
    exit 1
fi

echo "📥 [ ⟐ BISHOP_CORE ] Retrieving credentials from cloud residency..."
# Direct download logic for Google Drive
curl -L "https://docs.google.com/uc?export=download&id=${FILE_ID}" -o ${OUTPUT}

if [ -f "$OUTPUT" ]; then
    echo "✅ Credentials mirrored locally."
else
    echo "❌ Failed to retrieve credentials."
fi
