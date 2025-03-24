#!/bin/bash

# Navigate to project root
cd "$(dirname "$0")/.."

echo "===== Fixing go.sum for CLIBoard ====="

# Ensure go.mod is correct first
cat > go.mod << 'EOL'
module github.com/doko89/cliboard

go 1.21

require github.com/spf13/cobra v1.8.0

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)
EOL

# Remove existing go.sum if it exists
echo "1. Removing existing go.sum if present..."
rm -f go.sum

# Create a temporary main file that imports all required packages
echo "2. Creating temporary file with all imports..."
TMP_FILE=$(mktemp)
cat > $TMP_FILE << 'EOL'
// Temporary file to force dependency resolution
package main

import (
	"fmt"
	"github.com/spf13/cobra"
	_ "github.com/spf13/pflag"
	_ "github.com/inconshreveable/mousetrap"
)

func main() {
	fmt.Println("Resolving dependencies for", cobra.CompactErrorMessage)
}
EOL

# Download dependencies with verbose output
echo "3. Downloading dependencies..."
go mod download -x all

# Build the temporary file to generate go.sum
echo "4. Building temporary file to generate go.sum..."
go build -o /dev/null $TMP_FILE

# Clean up
rm -f $TMP_FILE

# Verify go.sum was created
if [ -f "go.sum" ]; then
    echo "✅ go.sum successfully created!"
    echo "Entries in go.sum: $(wc -l < go.sum)"
else
    echo "❌ Failed to create go.sum"
    exit 1
fi

echo "===== Fix complete! Try building the project now. ====="
