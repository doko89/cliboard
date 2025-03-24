#!/bin/bash

# Navigate to project root
cd "$(dirname "$0")/.."

echo "Fixing dependencies for CLIBoard..."

# Update module paths in code (if any are still using old paths)
echo "Updating import paths..."
find . -name "*.go" -type f -exec sed -i 's|github.com/doko/cliboard|github.com/doko89/cliboard|g' {} \;

# Update go.mod and download dependencies
echo "Updating go.mod and downloading dependencies..."
go mod tidy

echo "Dependencies fixed successfully. You can now build the project."
