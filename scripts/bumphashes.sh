#!/bin/bash

# Get the current directory
DIR="$(pwd)"

# Build the Go dependencies and capture the expected hash
GO_HASH=$(nix build .#panel 2>&1 | grep -o 'got: .*' | awk '{print $2}')

# Build the Node dependencies and capture the expected hash
NODE_HASH=$(nix build .#panelNodeModules 2>&1 | grep -o 'got: .*' | awk '{print $2}')

# Update the hashes.json file if at least one hash was captured
if [ -n "$GO_HASH" ] || [ -n "$NODE_HASH" ]; then
  # Read the current hashes
  CURRENT_HASHES=$(jq -c '.' ${DIR}/hashes.json)

  # Update the go hash if it was captured
  if [ -n "$GO_HASH" ]; then
    CURRENT_HASHES=$(echo $CURRENT_HASHES | jq --arg goHash "$GO_HASH" '.go = $goHash')
  fi

  # Update the node hash if it was captured
  if [ -n "$NODE_HASH" ]; then
    CURRENT_HASHES=$(echo $CURRENT_HASHES | jq --arg nodeHash "$NODE_HASH" '.node = $nodeHash')
  fi

  # Write the updated hashes to the file
  echo $CURRENT_HASHES > ${DIR}/hashes.json.tmp && mv ${DIR}/hashes.json.tmp ${DIR}/hashes.json

  echo "Updated hashes:"
  cat ${DIR}/hashes.json
else
  echo "Failed to capture both hashes."
fi
