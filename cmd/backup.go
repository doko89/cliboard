package cmd

import (
	"github.com/doko/cliboard/internal/backup"
	"github.com/spf13/cobra"
)

var enableBackupCmd = &cobra.Command{
	Use:   "enable-backup [domain]",
	Short: "Enable automatic site backup",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		return backup.EnableSite(domain)
	},
}

var disableBackupCmd = &cobra.Command{
	Use:   "disable-backup [domain]",
	Short: "Disable automatic site backup",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		return backup.DisableSite(domain)
	},
}

var enableDbBackupCmd = &cobra.Command{
	Use:   "enable-dbbackup",
	Short: "Enable automatic database backup",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return backup.EnableDatabase()
	},
}

var disableDbBackupCmd = &cobra.Command{
	Use:   "disable-dbbackup",
	Short: "Disable automatic database backup",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return backup.DisableDatabase()
	},
}
