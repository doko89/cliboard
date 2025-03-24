package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/doko89/cliboard/internal/config"
	"github.com/doko89/cliboard/internal/site"
	"github.com/spf13/cobra"
)

var createSiteCmd = &cobra.Command{
	Use:   "create-site [domain]",
	Short: "Create a new site",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		return site.Create(domain)
	},
}

var deleteSiteCmd = &cobra.Command{
	Use:   "delete-site [domain]",
	Short: "Delete an existing site",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		return site.Delete(domain)
	},
}

var webrootCmd = &cobra.Command{
	Use:   "webroot",
	Short: "Manage site webroot",
}

var webrootUpdateCmd = &cobra.Command{
	Use:   "update [domain] [path]",
	Short: "Update site webroot path",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		path := args[1]
		return site.UpdateWebroot(domain, path)
	},
}

func init() {
	webrootCmd.AddCommand(webrootUpdateCmd)
}
