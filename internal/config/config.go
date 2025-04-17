// internal/config/config.go
package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	APIPort              string        `mapstructure:"API_PORT"`
	JWTSecret            string        `mapstructure:"JWT_SECRET"`
	JWTExpirationMinutes time.Duration `mapstructure:"JWT_EXPIRATION_MINUTES"`
	APIUserGroup         string        `mapstructure:"API_USER_GROUP"`  // <-- NEW: Group required for basic API login (alternative to clab_admins)
	SuperuserGroup       string        `mapstructure:"SUPERUSER_GROUP"` // <-- REMAINS: Group for elevated privileges (view all labs, etc.)
	ClabRuntime          string        `mapstructure:"CLAB_RUNTIME"`
	LogLevel             string        `mapstructure:"LOG_LEVEL"`
	// --- New TLS Fields ---
	TLSEnable   bool   `mapstructure:"TLS_ENABLE"`
	TLSCertFile string `mapstructure:"TLS_CERT_FILE"`
	TLSKeyFile  string `mapstructure:"TLS_KEY_FILE"`
	// --- Gin Settings ---
	GinMode        string `mapstructure:"GIN_MODE"`
	TrustedProxies string `mapstructure:"TRUSTED_PROXIES"` // Comma-separated list or "nil" to disable
}

var AppConfig Config

func LoadConfig() error {
	viper.SetConfigFile(".env") // Look for .env file
	viper.AutomaticEnv()        // Read from environment variables as fallback/override

	// --- Set Defaults ---
	viper.SetDefault("API_PORT", "8080")
	viper.SetDefault("JWT_SECRET", "default_secret_change_me")
	viper.SetDefault("JWT_EXPIRATION_MINUTES", 60)
	viper.SetDefault("API_USER_GROUP", "")
	viper.SetDefault("SUPERUSER_GROUP", "")
	viper.SetDefault("CLAB_RUNTIME", "docker")
	viper.SetDefault("LOG_LEVEL", "info")
	// --- New TLS Defaults ---
	viper.SetDefault("TLS_ENABLE", false) // Disabled by default
	viper.SetDefault("TLS_CERT_FILE", "") // No default paths
	viper.SetDefault("TLS_KEY_FILE", "")
	// --- Gin Settings Defaults ---
	viper.SetDefault("GIN_MODE", "debug")   // Use 'release' for production
	viper.SetDefault("TRUSTED_PROXIES", "") // Empty means trust all, "nil" means trust none

	err := viper.ReadInConfig()
	// Ignore if .env file not found, rely on defaults/env vars
	if _, ok := err.(viper.ConfigFileNotFoundError); !ok && err != nil {
		return err
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		return err
	}

	// Convert minutes to duration
	AppConfig.JWTExpirationMinutes = AppConfig.JWTExpirationMinutes * time.Minute

	return nil
}
