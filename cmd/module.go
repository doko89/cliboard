package cmd

import (
	"github.com/doko/cliboard/internal/module"
	"github.com/spf13/cobra"
)

var addModuleCmd = &cobra.Command{
	Use:   "add-module [domain] [module]",
	Short: "Add a Caddy module to a site",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		moduleName := args[1]
		return module.Add(domain, moduleName)
	},
}

var removeModuleCmd = &cobra.Command{
	Use:   "remove-module [domain] [module]",
	Short: "Remove a Caddy module from a site",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		moduleName := args[1]
		return module.Remove(domain, moduleName)
	},
}

var listModulesCmd = &cobra.Command{
	Use:   "list-modules [domain]",
	Short: "List active modules for a site",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		return module.List(domain)
	},
}

var listAvailableModulesCmd = &cobra.Command{
	Use:   "list-available-modules",
	Short: "List all available Caddy modules",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return module.ListAvailable()
	},
}
