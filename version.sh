#!/bin/bash

# TEMP
# Script execute by Github action to increase the version
# in PR

# Determine the current version
VERSION=$( cat manifest.json | jq -r '.version')
IFS='.'

read -a strarr <<< "$VERSION"

MAJOR=$( echo "${strarr[0]}" | sed 's/^.//' )
MINOR=${strarr[1]}
PATCH=${strarr[2]}

# Find the label
for label in ${{ github.event.pull_request.labels }}
do
    echo "$label" | jq -r '.name'
done

# Build the new version

# Commit