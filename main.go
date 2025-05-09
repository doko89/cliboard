package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/doko89/cliboard/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
