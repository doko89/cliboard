package caddy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/doko89/cliboard/internal/config"
)

// Reload reloads the Caddy server
func Reload() error {
	// Check if Caddy is installed
	if _, err := exec.LookPath("caddy"); err != nil {
		return fmt.Errorf("Caddy is not installed")
	}

	// Run caddy reload
	cmd := exec.Command("systemctl", "reload", "caddy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reload Caddy: %v", err)
	}

	return nil
}

// Install installs Caddy with the necessary configuration
func Install() error {
	// Check if Caddy is already installed
	if _, err := exec.LookPath("caddy"); err == nil {
		fmt.Println("Caddy is already installed")
		return nil
	}

	fmt.Println("Installing Caddy...")

	// Add Caddy repository
	cmd := exec.Command("sh", "-c", `
    apt-get update
    apt-get install -y debian-keyring debian-archive-keyring apt-transport-https curl
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
    apt-get update
    apt-get install -y caddy
    `)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Caddy: %v", err)
	}

	// Create directory structure
	dirs := []string{
		config.CaddyRootDir,
		config.CaddyModulesDir,
		config.CaddyPHPDir,
		config.CaddySitesDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Create main Caddy configuration
	caddyConfig := `{
    admin off
    log {
        output file /var/log/caddy/access.log
        format json
    }
    email admin@localhost
}

(common) {
    log {
        output file /var/log/caddy/{host}.access.log
        format json
    }
    header ?Server "CLIBoard"
    encode gzip
}

import modules.d/*
import php.d/*
import sites.d/*
`

	if err := os.WriteFile(filepath.Join(config.CaddyRootDir, "Caddyfile"), []byte(caddyConfig), 0644); err != nil {
		return fmt.Errorf("failed to create Caddy configuration: %v", err)
	}

	// Create default modules
	defaultModules := map[string]string{
		"cache-headers": `(cache-headers) {
    header Cache-Control "public, max-age=3600"
}`,
		"compression": `(compression) {
    encode zstd gzip
}`,
		"local-access": `(local-access) {
    @local {
        remote_ip 127.0.0.1
        remote_ip 10.0.0.0/8
        remote_ip 172.16.0.0/12
        remote_ip 192.168.0.0/16
    }
}`,
		"ratelimit": `(ratelimit) {
    rate_limit {
        zone dynamic {
            key {remote_host}
            events 10
            window 10s
        }
    }
}`,
		"security": `(security) {
    header {
        X-Content-Type-Options "nosniff"
        X-Frame-Options "SAMEORIGIN"
        X-XSS-Protection "1; mode=block"
        Referrer-Policy "strict-origin-when-cross-origin"
    }
}`,
		"spa": `(spa) {
    try_files {path} /index.html
}`,
		"static_cache": `(static_cache) {
    @static {
        file {
            try_files {path}
        }
        path *.ico *.css *.js *.gif *.jpg *.jpeg *.png *.svg *.woff *.woff2
    }
    header @static Cache-Control "public, max-age=86400"
}`,
	}

	for name, content := range defaultModules {
		if err := os.WriteFile(filepath.Join(config.CaddyModulesDir, name), []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create module %s: %v", name, err)
		}
	}

	// Create sites directory
	if err := os.MkdirAll(config.SitesRootDir, 0755); err != nil {
		return fmt.Errorf("failed to create sites directory: %v", err)
	}

	// Restart Caddy
	cmd = exec.Command("systemctl", "restart", "caddy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart Caddy: %v", err)
	}

	fmt.Println("Caddy installed successfully")
	return nil
}
