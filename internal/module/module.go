package module

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/doko/cliboard/internal/caddy"
	"github.com/doko/cliboard/internal/config"
)

// Add adds a module to a site configuration
func Add(domain, moduleName string) error {
	// Check if module exists
	modulePath := config.GetModulePath(moduleName)
	if _, err := os.Stat(modulePath); os.IsNotExist(err) {
		return fmt.Errorf("module %s does not exist", moduleName)
	}

	// Check if site exists
	siteConfigPath := config.GetSiteConfigPath(domain)
	if _, err := os.Stat(siteConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("site %s does not exist", domain)
	}

	// Read site configuration
	siteConfig, err := os.ReadFile(siteConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read site configuration: %v", err)
	}

	// Check if module is already imported
	if strings.Contains(string(siteConfig), fmt.Sprintf("import %s", moduleName)) {
		return fmt.Errorf("module %s is already enabled for site %s", moduleName, domain)
	}

	// Add module import to site configuration
	updatedConfig := strings.Replace(string(siteConfig),
		fmt.Sprintf("%s {", domain),
		fmt.Sprintf("%s {\n    import %s", domain, moduleName),
		1)

	// Write updated configuration
	if err := os.WriteFile(siteConfigPath, []byte(updatedConfig), 0644); err != nil {
		return fmt.Errorf("failed to update site configuration: %v", err)
	}

	// Reload Caddy to apply changes
	if err := caddy.Reload(); err != nil {
		return fmt.Errorf("failed to reload Caddy: %v", err)
	}

	fmt.Printf("Module %s added to site %s successfully\n", moduleName, domain)
	return nil
}

// Remove removes a module from a site configuration
func Remove(domain, moduleName string) error {
	// Check if site exists
	siteConfigPath := config.GetSiteConfigPath(domain)
	if _, err := os.Stat(siteConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("site %s does not exist", domain)
	}

	// Read site configuration
	siteConfig, err := os.ReadFile(siteConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read site configuration: %v", err)
	}

	// Check if module is imported
	importLine := fmt.Sprintf("import %s", moduleName)
	if !strings.Contains(string(siteConfig), importLine) {
		return fmt.Errorf("module %s is not enabled for site %s", moduleName, domain)
	}

	// Remove module import from site configuration
	updatedConfig := strings.Replace(string(siteConfig),
		fmt.Sprintf("    import %s\n", moduleName),
		"",
		1)

	// Write updated configuration
	if err := os.WriteFile(siteConfigPath, []byte(updatedConfig), 0644); err != nil {
		return fmt.Errorf("failed to update site configuration: %v", err)
	}

	// Reload Caddy to apply changes
	if err := caddy.Reload(); err != nil {
		return fmt.Errorf("failed to reload Caddy: %v", err)
	}

	fmt.Printf("Module %s removed from site %s successfully\n", moduleName, domain)
	return nil
}

// List lists active modules for a site
func List(domain string) error {
	// Check if site exists
	siteConfigPath := config.GetSiteConfigPath(domain)
	if _, err := os.Stat(siteConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("site %s does not exist", domain)
	}

	// Read site configuration
	siteConfig, err := os.ReadFile(siteConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read site configuration: %v", err)
	}

	// Extract module imports
	var modules []string
	lines := strings.Split(string(siteConfig), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "import ") {
			moduleName := strings.TrimPrefix(line, "import ")
			modules = append(modules, moduleName)
		}
	}

	// Print modules
	if len(modules) == 0 {
		fmt.Printf("No active modules for site %s\n", domain)
	} else {
		fmt.Printf("Active modules for site %s:\n", domain)
		for _, module := range modules {
			fmt.Printf("- %s\n", module)
		}
	}

	return nil
}

// ListAvailable lists all available Caddy modules
func ListAvailable() error {
	files, err := ioutil.ReadDir(config.CaddyModulesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("modules directory does not exist: %s", config.CaddyModulesDir)
		}
		return fmt.Errorf("failed to read modules directory: %v", err)
	}

	if len(files) == 0 {
		fmt.Println("No available modules found")
		return nil
	}

	fmt.Println("Available Caddy modules:")
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fmt.Printf("- %s\n", file.Name())
	}

	return nil
}
