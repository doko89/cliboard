package config

const (
	// Site directories
	SitesRootDir = "/apps/sites"
	
	// Caddy configuration directories
	CaddyRootDir    = "/etc/caddy"
	CaddyModulesDir = "/etc/caddy/modules.d"
	CaddyPHPDir     = "/etc/caddy/php.d"
	CaddySitesDir   = "/etc/caddy/sites.d"
	
	// Backup directories
	BackupDailyDir  = "/backup/daily"
	BackupWeeklyDir = "/backup/weekly"
)

// GetSiteDirectory returns the full directory path for a site
func GetSiteDirectory(domain string) string {
	return SitesRootDir + "/" + domain
}

// GetSiteConfigPath returns the Caddy configuration file path for a site
func GetSiteConfigPath(domain string) string {
	return CaddySitesDir + "/" + domain + ".caddy"
}

// GetPHPConfigPath returns the PHP configuration file path for a specific PHP version
func GetPHPConfigPath(version string) string {
	return CaddyPHPDir + "/php" + version + "_config"
}

// GetModulePath returns the path to a Caddy module configuration
func GetModulePath(module string) string {
	return CaddyModulesDir + "/" + module
}

// GetBackupDailyPath returns the daily backup path for a domain
func GetBackupDailyPath(domain string) string {
	return BackupDailyDir + "/" + domain
}

// GetBackupWeeklyPath returns the weekly backup path for a domain
func GetBackupWeeklyPath(domain string) string {
	return BackupWeeklyDir + "/" + domain
}
