package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cliboard",
	Short: "CLIBoard - Web Control Panel for VPS Management",
	Long: `CLIBoard is a CLI-based web control panel that helps manage servers/VPS 
using Caddy as the web server. It supports site creation, PHP management, 
Caddy modules, automatic backups, and more.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add commands
	rootCmd.AddCommand(createSiteCmd)
	rootCmd.AddCommand(deleteSiteCmd)
	rootCmd.AddCommand(addModuleCmd)
	rootCmd.AddCommand(removeModuleCmd)
	rootCmd.AddCommand(listModulesCmd)
	rootCmd.AddCommand(listAvailableModulesCmd)
	rootCmd.AddCommand(webrootCmd)
	rootCmd.AddCommand(phpCmd)
	rootCmd.AddCommand(enableBackupCmd)
	rootCmd.AddCommand(disableBackupCmd)
	rootCmd.AddCommand(enableDbBackupCmd)
	rootCmd.AddCommand(disableDbBackupCmd)
	// Version and completion commands are added in their respective files
}
