#!/bin/bash

set -x

# Configure clob, perpetual and market module genesis

# Define variables
GENESIS_FILE="$HOME/.nemo/config/genesis.json"
TMP_GENESIS_FILE="./tmp-genesis-update.json"

# Load the updates JSON structure (if any other updates, add here)
UPDATES_FILE="./genesis-update.json"

clob=$(jq '.clob' $UPDATES_FILE)
perpetuals=$(jq '.perpetuals' $UPDATES_FILE)
prices=$(jq '.prices' $UPDATES_FILE)

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo "Error: jq is required but not installed. Please install jq."
    exit 1
fi

updated_json=$(jq ".app_state.clob = ${clob}"  $GENESIS_FILE)
updated_json=$(echo "$updated_json" | jq ".app_state.perpetuals = ${perpetuals}")
updated_json=$(echo "$updated_json" | jq ".app_state.prices = ${prices}")

# Write the updated JSON object to the temp file and move it to the original location
echo "$updated_json" > $TMP_GENESIS_FILE
mv $TMP_GENESIS_FILE $GENESIS_FILE

# Echo the updated JSON for debugging
echo "FINAL GENESIS:"
cat $GENESIS_FILE | jq '.'