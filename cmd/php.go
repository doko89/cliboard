package cmd

import (
	"github.com/doko89/cliboard/internal/php"
	"github.com/spf13/cobra"
)

var phpCmd = &cobra.Command{
	Use:   "php",
	Short: "Manage PHP versions and modules",
}

var enablePhpCmd = &cobra.Command{
	Use:   "enable-php [domain] [version]",
	Short: "Enable PHP for a site with specified version",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		version := args[1]
		return php.Enable(domain, version)
	},
}

var disablePhpCmd = &cobra.Command{
	Use:   "disable-php [domain]",
	Short: "Disable PHP for a site",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		return php.Disable(domain)
	},
}

var updatePhpCmd = &cobra.Command{
	Use:   "update-php [domain]",
	Short: "Update PHP configuration for a site",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		return php.Update(domain)
	},
}

var phpInstallCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install a specific PHP version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		return php.Install(version)
	},
}

var phpUninstallCmd = &cobra.Command{
	Use:   "uninstall [version]",
	Short: "Uninstall a specific PHP version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		return php.Uninstall(version)
	},
}

var phpListInstalledCmd = &cobra.Command{
	Use:   "list-installed",
	Short: "List installed PHP versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return php.ListInstalled()
	},
}

var phpModuleCmd = &cobra.Command{
	Use:   "module",
	Short: "Manage PHP modules",
}

var phpModuleListAvailableCmd = &cobra.Command{
	Use:   "list-available [version]",
	Short: "List available modules for a PHP version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		return php.ListAvailableModules(version)
	},
}

var phpModuleAddCmd = &cobra.Command{
	Use:   "add [version] [module]",
	Short: "Add a module to a PHP version",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		module := args[1]
		return php.AddModule(version, module)
	},
}

var phpModuleRemoveCmd = &cobra.Command{
	Use:   "remove [version] [module]",
	Short: "Remove a module from a PHP version",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		module := args[1]
		return php.RemoveModule(version, module)
	},
}

func init() {
	rootCmd.AddCommand(enablePhpCmd)
	rootCmd.AddCommand(disablePhpCmd)
	rootCmd.AddCommand(updatePhpCmd)
	
	phpCmd.AddCommand(phpInstallCmd)
	phpCmd.AddCommand(phpUninstallCmd)
	phpCmd.AddCommand(phpListInstalledCmd)
	phpCmd.AddCommand(phpModuleCmd)
	
	phpModuleCmd.AddCommand(phpModuleListAvailableCmd)
	phpModuleCmd.AddCommand(phpModuleAddCmd)
	phpModuleCmd.AddCommand(phpModuleRemoveCmd)
}
