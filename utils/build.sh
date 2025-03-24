#!/bin/bash

# Navigate to project root
cd "$(dirname "$0")/.."

# Make sure dependencies are fixed
echo "Fixing dependencies..."
./utils/fix_go_mod.sh

# Set variables
VERSION=${VERSION:-"1.0.0"}
GOOS=${GOOS:-"linux"}
GOARCH=${GOARCH:-"amd64"}
BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
COMMIT_SHA=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo "Building CLIBoard v${VERSION} for ${GOOS}/${GOARCH}..."
mkdir -p dist

# Build the binary
go build -o dist/cliboard-${GOOS}-${GOARCH} \
  -ldflags "-X github.com/doko89/cliboard/cmd.Version=${VERSION} -X github.com/doko89/cliboard/cmd.BuildDate=${BUILD_DATE} -X github.com/doko89/cliboard/cmd.CommitSHA=${COMMIT_SHA}" \
  main.go

if [ $? -eq 0 ]; then
  echo "Build successful. Binary located at dist/cliboard-${GOOS}-${GOARCH}"
else
  echo "Build failed."
  exit 1
fi
