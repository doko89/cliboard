#!/bin/bash

# Navigate to project root
cd "$(dirname "$0")/.."

echo "===== Force fixing dependency issues ====="

# 1. Remove go.sum and go.mod
echo "1. Removing existing go.mod and go.sum files..."
rm -f go.sum go.mod

# 2. Create a proper go.mod file
echo "2. Creating new go.mod file..."
cat > go.mod << 'EOL'
module github.com/doko89/cliboard

go 1.21

require (
	github.com/spf13/cobra v1.8.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)
EOL

# 3. Initialize a completely new module
echo "3. Forcing module initialization..."
go mod download

# 4. Create a temporary file that forces importing all dependencies
echo "4. Creating temporary file to test imports..."
TMP_DIR=$(mktemp -d)
cat > $TMP_DIR/main.go << 'EOL'
package main

import (
	"fmt"
	
	"github.com/spf13/cobra"
	_ "github.com/spf13/pflag"
	_ "github.com/inconshreveable/mousetrap"
)

func main() {
	fmt.Println("Testing imports with", cobra.CompactErrorMessage)
}
EOL

# 5. Build the temporary file to force dependency resolution
echo "5. Building temporary file to generate go.sum..."
cd $TMP_DIR
go mod init temp
go mod edit -replace github.com/doko89/cliboard=../../
go mod tidy
go build -o /dev/null ./main.go
cd - > /dev/null

# 6. Run go mod tidy in our project
echo "6. Running go mod tidy in our project..."
go mod tidy

# 7. Verify go.sum was created
if [ -f "go.sum" ]; then
    echo "✅ go.sum successfully created!"
    echo "Entries in go.sum: $(wc -l < go.sum)"
else
    echo "❌ Failed to create go.sum"
    exit 1
fi

# 8. Clean up
rm -rf $TMP_DIR

echo "===== Dependency fix complete! Try building the project now. ====="
