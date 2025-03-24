package php

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/doko/cliboard/internal/caddy"
	"github.com/doko/cliboard/internal/config"
)

// Enable enables PHP for a site with the specified version
func Enable(domain, version string) error {
	// Check if site exists
	siteConfigPath := config.GetSiteConfigPath(domain)
	if _, err := os.Stat(siteConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("site %s does not exist", domain)
	}

	// Check if PHP version is installed
	if !isVersionInstalled(version) {
		return fmt.Errorf("PHP version %s is not installed", version)
	}

	// Check if PHP configuration exists
	phpConfigName := fmt.Sprintf("php%s_config", version)
	phpConfigPath := config.GetPHPConfigPath(version)
	if _, err := os.Stat(phpConfigPath); os.IsNotExist(err) {
		// Create PHP configuration if it doesn't exist
		phpConfig := fmt.Sprintf(`(php%s_config) {
    php_fastcgi unix//run/php/php%s-fpm.sock
}
`, version, version)
		if err := os.MkdirAll(config.CaddyPHPDir, 0755); err != nil {
			return fmt.Errorf("failed to create PHP configuration directory: %v", err)
		}
		if err := os.WriteFile(phpConfigPath, []byte(phpConfig), 0644); err != nil {
			return fmt.Errorf("failed to create PHP configuration: %v", err)
		}
	}

	// Read site configuration
	siteConfig, err := os.ReadFile(siteConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read site configuration: %v", err)
	}

	// Check if any PHP version is already enabled
	for _, ver := range getInstalledVersions() {
		phpImport := fmt.Sprintf("import php%s_config", ver)
		if strings.Contains(string(siteConfig), phpImport) {
			// Replace the existing PHP version with the new one
			updatedConfig := strings.Replace(string(siteConfig),
				phpImport,
				fmt.Sprintf("import %s", phpConfigName),
				1)
			
			if err := os.WriteFile(siteConfigPath, []byte(updatedConfig), 0644); err != nil {
				return fmt.Errorf("failed to update site configuration: %v", err)
			}
			
			// Reload Caddy
			if err := caddy.Reload(); err != nil {
				return fmt.Errorf("failed to reload Caddy: %v", err)
			}
			
			fmt.Printf("PHP version updated to %s for site %s\n", version, domain)
			return nil
		}
	}

	// PHP not yet enabled, add it to site configuration
	updatedConfig := strings.Replace(string(siteConfig),
		fmt.Sprintf("%s {", domain),
		fmt.Sprintf("%s {\n    import %s", domain, phpConfigName),
		1)

	if err := os.WriteFile(siteConfigPath, []byte(updatedConfig), 0644); err != nil {
		return fmt.Errorf("failed to update site configuration: %v", err)
	}

	// Reload Caddy
	if err := caddy.Reload(); err != nil {
		return fmt.Errorf("failed to reload Caddy: %v", err)
	}

	fmt.Printf("PHP version %s enabled for site %s\n", version, domain)
	return nil
}

// Disable disables PHP for a site
func Disable(domain string) error {
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

	// Check if any PHP version is enabled
	phpEnabled := false
	updatedConfig := string(siteConfig)
	
	for _, ver := range getInstalledVersions() {
		phpImport := fmt.Sprintf("import php%s_config", ver)
		if strings.Contains(updatedConfig, phpImport) {
			updatedConfig = strings.Replace(updatedConfig,
				fmt.Sprintf("    %s\n", phpImport),
				"",
				1)
			phpEnabled = true
		}
	}

	if !phpEnabled {
		return fmt.Errorf("PHP is not enabled for site %s", domain)
	}

	if err := os.WriteFile(siteConfigPath, []byte(updatedConfig), 0644); err != nil {
		return fmt.Errorf("failed to update site configuration: %v", err)
	}

	// Reload Caddy
	if err := caddy.Reload(); err != nil {
		return fmt.Errorf("failed to reload Caddy: %v", err)
	}

	fmt.Printf("PHP disabled for site %s\n", domain)
	return nil
}

// Update updates PHP configuration for a site
func Update(domain string) error {
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

	// Check which PHP version is enabled
	var currentVersion string
	for _, ver := range getInstalledVersions() {
		phpImport := fmt.Sprintf("import php%s_config", ver)
		if strings.Contains(string(siteConfig), phpImport) {
			currentVersion = ver
			break
		}
	}

	if currentVersion == "" {
		return fmt.Errorf("PHP is not enabled for site %s", domain)
	}

	// Re-enable the current version (will recreate the PHP configuration if needed)
	return Enable(domain, currentVersion)
}

// Install installs a specific PHP version
func Install(version string) error {
	// Check if already installed
	if isVersionInstalled(version) {
		return fmt.Errorf("PHP version %s is already installed", version)
	}
	
	fmt.Printf("Installing PHP %s...\n", version)
	
	// Install PHP packages
	cmd := exec.Command("apt-get", "install", "-y", 
		fmt.Sprintf("php%s", version),
		fmt.Sprintf("php%s-fpm", version),
		fmt.Sprintf("php%s-common", version),
		fmt.Sprintf("php%s-cli", version))
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install PHP %s: %v", version, err)
	}
	
	// Enable the PHP-FPM service
	cmd = exec.Command("systemctl", "enable", fmt.Sprintf("php%s-fpm", version))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to enable PHP-FPM service: %v", err)
	}
	
	// Start the PHP-FPM service
	cmd = exec.Command("systemctl", "start", fmt.Sprintf("php%s-fpm", version))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start PHP-FPM service: %v", err)
	}
	
	// Create PHP configuration for Caddy
	phpConfig := fmt.Sprintf(`(php%s_config) {
    php_fastcgi unix//run/php/php%s-fpm.sock
}
`, version, version)

	if err := os.MkdirAll(config.CaddyPHPDir, 0755); err != nil {
		return fmt.Errorf("failed to create PHP configuration directory: %v", err)
	}

	phpConfigPath := config.GetPHPConfigPath(version)
	if err := os.WriteFile(phpConfigPath, []byte(phpConfig), 0644); err != nil {
		return fmt.Errorf("failed to create PHP configuration: %v", err)
	}
	
	fmt.Printf("PHP %s installed successfully\n", version)
	return nil
}

// Uninstall uninstalls a specific PHP version
func Uninstall(version string) error {
	// Check if installed
	if !isVersionInstalled(version) {
		return fmt.Errorf("PHP version %s is not installed", version)
	}
	
	fmt.Printf("Uninstalling PHP %s...\n", version)
	
	// Stop and disable the PHP-FPM service
	exec.Command("systemctl", "stop", fmt.Sprintf("php%s-fpm", version)).Run()
	exec.Command("systemctl", "disable", fmt.Sprintf("php%s-fpm", version)).Run()
	
	// Remove PHP packages
	cmd := exec.Command("apt-get", "remove", "-y", 
		fmt.Sprintf("php%s", version),
		fmt.Sprintf("php%s-fpm", version),
		fmt.Sprintf("php%s-common", version))
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to uninstall PHP %s: %v", version, err)
	}
	
	// Remove PHP configuration for Caddy
	phpConfigPath := config.GetPHPConfigPath(version)
	if err := os.Remove(phpConfigPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove PHP configuration: %v", err)
	}
	
	fmt.Printf("PHP %s uninstalled successfully\n", version)
	return nil
}

// ListInstalled lists installed PHP versions
func ListInstalled() error {
	versions := getInstalledVersions()
	
	if len(versions) == 0 {
		fmt.Println("No PHP versions installed")
		return nil
	}
	
	fmt.Println("Installed PHP versions:")
	for _, version := range versions {
		fmt.Printf("- %s\n", version)
	}
	
	return nil
}

// ListAvailableModules lists available modules for a PHP version
func ListAvailableModules(version string) error {
	// Check if PHP version is installed
	if !isVersionInstalled(version) {
		return fmt.Errorf("PHP version %s is not installed", version)
	}
	
	cmd := exec.Command("apt-cache", "search", fmt.Sprintf("php%s-", version))
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list available modules: %v", err)
	}
	
	modules := []string{}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, fmt.Sprintf("php%s-", version)) {
			parts := strings.SplitN(line, " - ", 2)
			if len(parts) > 0 {
				moduleName := strings.TrimPrefix(parts[0], fmt.Sprintf("php%s-", version))
				modules = append(modules, moduleName)
			}
		}
	}
	
	if len(modules) == 0 {
		fmt.Printf("No available modules found for PHP %s\n", version)
		return nil
	}
	
	fmt.Printf("Available modules for PHP %s:\n", version)
	for _, module := range modules {
		fmt.Printf("- %s\n", module)
	}
	
	return nil
}

// AddModule adds a module to a PHP version
func AddModule(version, module string) error {
	// Check if PHP version is installed
	if !isVersionInstalled(version) {
		return fmt.Errorf("PHP version %s is not installed", version)
	}
	
	// Install PHP module
	packageName := fmt.Sprintf("php%s-%s", version, module)
	fmt.Printf("Installing PHP module %s...\n", packageName)
	
	cmd := exec.Command("apt-get", "install", "-y", packageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install PHP module %s: %v", module, err)
	}
	
	// Restart PHP-FPM
	cmd = exec.Command("systemctl", "restart", fmt.Sprintf("php%s-fpm", version))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart PHP-FPM: %v", err)
	}
	
	fmt.Printf("PHP module %s installed successfully for PHP %s\n", module, version)
	return nil
}

// RemoveModule removes a module from a PHP version
func RemoveModule(version, module string) error {
	// Check if PHP version is installed
	if !isVersionInstalled(version) {
		return fmt.Errorf("PHP version %s is not installed", version)
	}
	
	// Remove PHP module
	packageName := fmt.Sprintf("php%s-%s", version, module)
	fmt.Printf("Removing PHP module %s...\n", packageName)
	
	cmd := exec.Command("apt-get", "remove", "-y", packageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove PHP module %s: %v", module, err)
	}
	
	// Restart PHP-FPM
	cmd = exec.Command("systemctl", "restart", fmt.Sprintf("php%s-fpm", version))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart PHP-FPM: %v", err)
	}
	
	fmt.Printf("PHP module %s removed successfully from PHP %s\n", module, version)
	return nil
}

// Helper functions

// isVersionInstalled checks if a PHP version is installed
func isVersionInstalled(version string) bool {
	cmd := exec.Command("which", fmt.Sprintf("php%s", version))
	if err := cmd.Run(); err == nil {
		return true
	}
	return false
}

// getInstalledVersions returns a list of installed PHP versions
func getInstalledVersions() []string {
	var versions []string
	
	// Check common PHP versions
	for _, ver := range []string{"7.0", "7.1", "7.2", "7.3", "7.4", "8.0", "8.1", "8.2", "8.3"} {
		if isVersionInstalled(ver) {
			versions = append(versions, ver)
		}
	}
	
	return versions
}
