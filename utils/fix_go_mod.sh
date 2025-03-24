#!/bin/bash

# Navigate to project root
cd "$(dirname "$0")/.."

echo "Creating temporary main file to force dependency resolution..."

# Create a temporary file to force Go to download all dependencies
cat > cmd/temp_dependency_resolver.go << 'EOL'
// Temporary file to force dependency resolution
// +build ignore

package cmd

import (
    "fmt"
    _ "github.com/spf13/cobra"
)

func DummyFunction() {
    fmt.Println("This function is just to force Go to resolve dependencies")
}
EOL

# Run go get to download dependencies
echo "Downloading dependencies..."
go get -v github.com/spf13/cobra@latest
go mod tidy

# Remove the temporary file
rm -f cmd/temp_dependency_resolver.go

echo "Dependencies fixed. Try building the project now."
