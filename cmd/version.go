package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information
var (
	Version   = "dev"
	BuildDate = "unknown"
	CommitSHA = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of CLIBoard",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("CLIBoard v%s\n", Version)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Commit: %s\n", CommitSHA)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
