#!/bin/bash

# Example:
# sh hello.sh http://www.go-minitwit.duckdns.org/version 2

BUILD_NUMBER_URL=$1
RELEASE_NUMEBER=$2

build_number() {
  CURL_RESULT=$(curl -s $BUILD_NUMBER_URL)
  CURRENT_NUMBER=$(echo $CURL_RESULT | jq ".")
  if [ "$CURRENT_NUMBER" -eq "$CURRENT_NUMBER" ] 2>/dev/null; then
    echo "Build number: $CURRENT_NUMBER ($(date))"
  else
    echo "Failed to get build number from url $BUILD_NUMBER_URL";
    exit 1
  fi
}

build_number
echo "Release build number: $RELEASE_NUMEBER"

END=$(($(date +%s)+60))
echo "Starting to wait for release: "$(date)
echo "Waiting for release until   : "$(date -r $END)


while true; do
  sleep 2
  build_number
  if [ "$CURRENT_NUMBER" -eq "$RELEASE_NUMEBER" ]; then
    echo "New release detected with build number $CURRENT_NUMBER"
    exit 0
  fi
  NOW=$(date +%s)
  if (("$NOW" > "$END")); then
    echo "Timeout waiting for new release"
    exit 1
  fi
done
