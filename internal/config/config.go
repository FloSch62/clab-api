// internal/config/config.go
package config

import (
	"fmt" // Import fmt
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	APIPort              string        `mapstructure:"API_PORT"`
	JWTSecret            string        `mapstructure:"JWT_SECRET"`
	JWTExpirationMinutes time.Duration `mapstructure:"JWT_EXPIRATION_MINUTES"`
	APIUserGroup         string        `mapstructure:"API_USER_GROUP"`  // Group required for basic API login (alternative to clab_admins)
	SuperuserGroup       string        `mapstructure:"SUPERUSER_GROUP"` // Group for elevated privileges
	ClabRuntime          string        `mapstructure:"CLAB_RUNTIME"`
	LogLevel             string        `mapstructure:"LOG_LEVEL"`
	TLSEnable            bool          `mapstructure:"TLS_ENABLE"`
	TLSCertFile          string        `mapstructure:"TLS_CERT_FILE"`
	TLSKeyFile           string        `mapstructure:"TLS_KEY_FILE"`
	GinMode              string        `mapstructure:"GIN_MODE"`
	TrustedProxies       string        `mapstructure:"TRUSTED_PROXIES"`
}

var AppConfig Config

// LoadConfig loads configuration from the specified .env file path and environment variables.
func LoadConfig(envFilePath string) error {
	// Use the provided file path
	viper.SetConfigFile(envFilePath)
	viper.AutomaticEnv() // Read from environment variables as fallback/override

	// --- Set Defaults ---
	viper.SetDefault("API_PORT", "8080")
	viper.SetDefault("JWT_SECRET", "default_secret_change_me")
	viper.SetDefault("JWT_EXPIRATION_MINUTES", 60)
	viper.SetDefault("API_USER_GROUP", "")
	viper.SetDefault("SUPERUSER_GROUP", "")
	viper.SetDefault("CLAB_RUNTIME", "docker")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("TLS_ENABLE", false)
	viper.SetDefault("TLS_CERT_FILE", "")
	viper.SetDefault("TLS_KEY_FILE", "")
	viper.SetDefault("GIN_MODE", "debug")
	viper.SetDefault("TRUSTED_PROXIES", "")

	err := viper.ReadInConfig()

	// Handle file not found error specifically
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// If the default file path was used and it wasn't found, it's okay.
			// If a specific path was provided via flag and it wasn't found, it's an error.
			defaultPath := ".env" // The default path we'll set in main.go
			if envFilePath != defaultPath {
				return fmt.Errorf("specified config file '%s' not found: %w", envFilePath, err)
			}
			// Otherwise (default path not found), just log a debug message and continue
			// (relying on env vars and defaults)
			// We'll add the logger initialization *after* config loading in main.go,
			// so we can't log here yet. We'll log the outcome in main.go.
		} else {
			// Some other error occurred reading the config file
			return fmt.Errorf("error reading config file '%s': %w", envFilePath, err)
		}
	}
	// If err is nil, the file was read successfully.

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		return fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Convert minutes to duration
	AppConfig.JWTExpirationMinutes = AppConfig.JWTExpirationMinutes * time.Minute

	return nil
}
