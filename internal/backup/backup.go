package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/doko/cliboard/internal/config"
)

// EnableSite enables automatic backup for a site
func EnableSite(domain string) error {
	// Check if site exists
	siteDir := config.GetSiteDirectory(domain)
	if _, err := os.Stat(siteDir); os.IsNotExist(err) {
		return fmt.Errorf("site %s does not exist", domain)
	}

	// Create backup directories
	dailyBackupDir := config.GetBackupDailyPath(domain)
	weeklyBackupDir := config.GetBackupWeeklyPath(domain)

	if err := os.MkdirAll(dailyBackupDir, 0755); err != nil {
		return fmt.Errorf("failed to create daily backup directory: %v", err)
	}

	if err := os.MkdirAll(weeklyBackupDir, 0755); err != nil {
		return fmt.Errorf("failed to create weekly backup directory: %v", err)
	}

	// Create cron jobs for daily and weekly backups
	dailyCron := fmt.Sprintf("0 1 * * * root rsync -a --delete --link-dest=%s/latest %s %s/$(date +%%Y%%m%%d) && ln -sf %s/$(date +%%Y%%m%%d) %s/latest\n",
		dailyBackupDir, siteDir, dailyBackupDir, dailyBackupDir, dailyBackupDir)

	weeklyCron := fmt.Sprintf("0 2 * * 0 root rsync -a --delete %s %s/$(date +%%Y%%m%%d) && ln -sf %s/$(date +%%Y%%m%%d) %s/latest\n",
		siteDir, weeklyBackupDir, weeklyBackupDir, weeklyBackupDir)

	// Write cron jobs to /etc/cron.d/
	cronFile := fmt.Sprintf("/etc/cron.d/cliboard-backup-%s", domain)
	
	cronContent := fmt.Sprintf("# CLIBoard backup cron jobs for %s\n%s%s", domain, dailyCron, weeklyCron)
	
	if err := os.WriteFile(cronFile, []byte(cronContent), 0644); err != nil {
		return fmt.Errorf("failed to create backup cron jobs: %v", err)
	}

	fmt.Printf("Automatic backups enabled for site %s\n", domain)
	return nil
}

// DisableSite disables automatic backup for a site
func DisableSite(domain string) error {
	// Remove cron jobs
	cronFile := fmt.Sprintf("/etc/cron.d/cliboard-backup-%s", domain)
	if err := os.Remove(cronFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove backup cron jobs: %v", err)
	}

	fmt.Printf("Automatic backups disabled for site %s\n", domain)
	return nil
}

// EnableDatabase enables automatic database backup
func EnableDatabase() error {
	// Check if MariaDB/MySQL is installed
	if !isDatabaseInstalled() {
		return fmt.Errorf("MariaDB/MySQL is not installed")
	}

	// Create backup directories
	dailyBackupDir := "/backup/daily/database"
	weeklyBackupDir := "/backup/weekly/database"

	if err := os.MkdirAll(dailyBackupDir, 0755); err != nil {
		return fmt.Errorf("failed to create daily backup directory: %v", err)
	}

	if err := os.MkdirAll(weeklyBackupDir, 0755); err != nil {
		return fmt.Errorf("failed to create weekly backup directory: %v", err)
	}

	// Create backup script
	backupScript := `#!/bin/bash

BACKUP_DIR="$1"
DATE=$(date +%Y%m%d)
MYSQL_USER="root"
MYSQL_PASSWORD=$(cat /root/.mysql_password 2>/dev/null || echo "")

# Create backup directory for today
mkdir -p "$BACKUP_DIR/$DATE"

# Get list of databases
DATABASES=$(mysql -u$MYSQL_USER ${MYSQL_PASSWORD:+-p$MYSQL_PASSWORD} -e "SHOW DATABASES;" | grep -Ev "(Database|information_schema|performance_schema)")

# Backup each database
for DB in $DATABASES; do
    mysqldump -u$MYSQL_USER ${MYSQL_PASSWORD:+-p$MYSQL_PASSWORD} --single-transaction --skip-lock-tables "$DB" > "$BACKUP_DIR/$DATE/$DB.sql"
done

# Create latest symlink
ln -sf "$BACKUP_DIR/$DATE" "$BACKUP_DIR/latest"
`

	backupScriptPath := "/usr/local/bin/cliboard-db-backup"
	if err := os.WriteFile(backupScriptPath, []byte(backupScript), 0755); err != nil {
		return fmt.Errorf("failed to create backup script: %v", err)
	}

	// Create cron jobs for daily and weekly backups
	dailyCron := fmt.Sprintf("0 3 * * * root /usr/local/bin/cliboard-db-backup %s\n", dailyBackupDir)
	weeklyCron := fmt.Sprintf("0 4 * * 0 root /usr/local/bin/cliboard-db-backup %s\n", weeklyBackupDir)

	// Write cron jobs to /etc/cron.d/
	cronFile := "/etc/cron.d/cliboard-db-backup"
	
	cronContent := fmt.Sprintf("# CLIBoard database backup cron jobs\n%s%s", dailyCron, weeklyCron)
	
	if err := os.WriteFile(cronFile, []byte(cronContent), 0644); err != nil {
		return fmt.Errorf("failed to create database backup cron jobs: %v", err)
	}

	fmt.Println("Automatic database backups enabled")
	return nil
}

// DisableDatabase disables automatic database backup
func DisableDatabase() error {
	// Remove cron jobs
	cronFile := "/etc/cron.d/cliboard-db-backup"
	if err := os.Remove(cronFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove database backup cron jobs: %v", err)
	}

	// Remove backup script
	backupScriptPath := "/usr/local/bin/cliboard-db-backup"
	if err := os.Remove(backupScriptPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove backup script: %v", err)
	}

	fmt.Println("Automatic database backups disabled")
	return nil
}

// Helper function to check if database is installed
func isDatabaseInstalled() bool {
	// Check for MariaDB
	if _, err := exec.LookPath("mariadb"); err == nil {
		return true
	}
	
	// Check for MySQL
	if _, err := exec.LookPath("mysql"); err == nil {
		return true
	}
	
	return false
}
