package site

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/doko89/cliboard/internal/caddy"
	"github.com/doko89/cliboard/internal/config"
	"github.com/doko89/cliboard/internal/utils"
)

// Create creates a new site with the given domain
func Create(domain string) error {
	// Create site directory
	siteDir := config.GetSiteDirectory(domain)
	if err := os.MkdirAll(siteDir, 0755); err != nil {
		return fmt.Errorf("failed to create site directory: %v", err)
	}

	// Create a sample index.html
	indexPath := filepath.Join(siteDir, "index.html")
	indexContent := fmt.Sprintf("<html><body><h1>Welcome to %s</h1><p>Site created with CLIBoard</p></body></html>", domain)
	if err := os.WriteFile(indexPath, []byte(indexContent), 0644); err != nil {
		return fmt.Errorf("failed to create index.html: %v", err)
	}

	// Create Caddy configuration
	caddyConfig := fmt.Sprintf(`%s {
    root * %s
    file_server
}
`, domain, siteDir)

	configPath := config.GetSiteConfigPath(domain)
	if err := os.WriteFile(configPath, []byte(caddyConfig), 0644); err != nil {
		return fmt.Errorf("failed to create site configuration: %v", err)
	}

	// Reload Caddy to apply changes
	if err := caddy.Reload(); err != nil {
		return fmt.Errorf("failed to reload Caddy: %v", err)
	}

	fmt.Printf("Site %s created successfully\n", domain)
	return nil
}

// Delete removes an existing site
func Delete(domain string) error {
	// Check if site exists
	siteDir := config.GetSiteDirectory(domain)
	if _, err := os.Stat(siteDir); os.IsNotExist(err) {
		return fmt.Errorf("site %s does not exist", domain)
	}

	// Ask for confirmation
	confirmed := utils.AskForConfirmation(fmt.Sprintf("Are you sure you want to delete site %s?", domain))
	if !confirmed {
		fmt.Println("Site deletion cancelled")
		return nil
	}

	// Remove site directory
	if err := os.RemoveAll(siteDir); err != nil {
		return fmt.Errorf("failed to remove site directory: %v", err)
	}

	// Remove Caddy configuration
	configPath := config.GetSiteConfigPath(domain)
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove site configuration: %v", err)
	}

	// Reload Caddy to apply changes
	if err := caddy.Reload(); err != nil {
		return fmt.Errorf("failed to reload Caddy: %v", err)
	}

	fmt.Printf("Site %s deleted successfully\n", domain)
	return nil
}

// UpdateWebroot updates the webroot path for a site
func UpdateWebroot(domain, path string) error {
	// Check if site exists
	siteDir := config.GetSiteDirectory(domain)
	if _, err := os.Stat(siteDir); os.IsNotExist(err) {
		return fmt.Errorf("site %s does not exist", domain)
	}

	// Make sure path starts with a slash
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Create the new webroot path if it doesn't exist
	newWebroot := filepath.Join(siteDir, strings.TrimPrefix(path, "/"))
	if err := os.MkdirAll(newWebroot, 0755); err != nil {
		return fmt.Errorf("failed to create webroot directory: %v", err)
	}

	// Update Caddy configuration
	configPath := config.GetSiteConfigPath(domain)
	config, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read site configuration: %v", err)
	}

	// Update the root directive
	updatedConfig := strings.Replace(string(config), 
		fmt.Sprintf("root * %s", siteDir), 
		fmt.Sprintf("root * %s", newWebroot), 
		1)

	if err := os.WriteFile(configPath, []byte(updatedConfig), 0644); err != nil {
		return fmt.Errorf("failed to update site configuration: %v", err)
	}

	// Reload Caddy to apply changes
	if err := caddy.Reload(); err != nil {
		return fmt.Errorf("failed to reload Caddy: %v", err)
	}

	fmt.Printf("Webroot for site %s updated to %s successfully\n", domain, newWebroot)
	return nil
}
