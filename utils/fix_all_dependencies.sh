#!/bin/bash

# Navigate to project root
cd "$(dirname "$0")/.."

echo "==== Comprehensive dependency fix for CLIBoard ===="

# Step 1: Update import paths in all Go files
echo "1. Updating all import paths..."
find . -name "*.go" -type f -exec sed -i 's|github.com/doko/cliboard|github.com/doko89/cliboard|g' {} \;

# Step 2: Clean module cache
echo "2. Cleaning module cache..."
go clean -modcache

# Step 3: Update go.mod file directly to ensure correct module path
echo "3. Updating go.mod file..."
cat > go.mod << 'EOL'
module github.com/doko89/cliboard

go 1.21

require github.com/spf13/cobra v1.8.0

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)
EOL

# Step 4: Explicitly download dependencies
echo "4. Downloading dependencies..."
go get -v github.com/spf13/cobra@v1.8.0
go get -v github.com/spf13/pflag@v1.0.5
go get -v github.com/inconshreveable/mousetrap@v1.1.0

# Step 5: Run go mod tidy to clean up
echo "5. Running go mod tidy..."
go mod tidy

# Step 6: Verify everything works by building a simple test
echo "6. Verifying build process..."
go build -o /dev/null main.go
if [ $? -eq 0 ]; then
    echo "✅ Success! Dependencies fixed and verified."
else
    echo "❌ Error: There are still dependency issues."
    exit 1
fi

echo "==== Fix complete! You can now build the project. ===="
