#!/usr/bin/env bash

COMMIT_DATE=$(git log -1 --pretty=format:"%ci" | awk '{print $1}')
COMMIT_TIME=$(git log -1 --pretty=format:"%ci" | awk '{print $2}' | sed 's/://g')
COMMIT_ID=$(git log -1 --pretty=format:"%h" | awk '{print $1}')

echo "${COMMIT_DATE}-${COMMIT_TIME}-${COMMIT_ID}"