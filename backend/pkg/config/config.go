package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/thienel/tlog"
)

// ServerConfig holds server configuration
type ServerConfig struct {
	Port        string `env:"PORT"`
	Env         string `env:"ENV"`
	ServiceName string `env:"SERVICE_NAME"`
	Version     string `env:"SERVICE_VERSION"`
}

// DatabaseConfig holds PostgreSQL configuration
type DatabaseConfig struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	DBName   string `env:"DB_NAME"`
	SSLMode  string `env:"DB_SSLMODE"`
	TimeZone string `env:"DB_TIMEZONE"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-default:"localhost"`
	Port     int    `env:"REDIS_PORT" env-default:"6379"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" env-default:"0"`
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret              string `env:"JWT_SECRET"`
	AccessExpiryMinutes int    `env:"JWT_ACCESS_EXPIRY_MINUTES"`
	RefreshExpiryHours  int    `env:"JWT_REFRESH_EXPIRY_HOURS"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level         string `env:"LOG_LEVEL"`
	EnableConsole bool   `env:"LOG_ENABLE_CONSOLE"`
	FilePath      string `env:"LOG_FILE_PATH"`
	MaxSizeMB     int    `env:"LOG_MAX_SIZE_MB"`
	MaxBackups    int    `env:"LOG_MAX_BACKUPS"`
	MaxAgeDays    int    `env:"LOG_MAX_AGE_DAYS"`
	Compress      bool   `env:"LOG_COMPRESS"`
}

// CookieConfig holds cookie configuration
type CookieConfig struct {
	Name        string `env:"COOKIE_NAME"`
	RefreshName string `env:"COOKIE_REFRESH_NAME"`
	Domain      string `env:"COOKIE_DOMAIN"`
	Secure      bool   `env:"COOKIE_SECURE"`
	SameSite    string `env:"COOKIE_SAMESITE"`
	Path        string `env:"COOKIE_PATH"`
}

// CasbinConfig holds Casbin authorization configuration
type CasbinConfig struct {
	ModelPath string `env:"CASBIN_MODEL_PATH"`
}

// NanoCMDConfig holds NanoCMD server configuration
type NanoCMDConfig struct {
	BaseURL  string `env:"NANOCMD_URL"`
	Username string `env:"NANOCMD_USERNAME"`
	Password string `env:"NANOCMD_PASSWORD"`
}

// NanoMDMConfig holds NanoMDM and NanoDEP server configuration
type NanoMDMConfig struct {
	MDMBaseURL         string `env:"NANOMDM_URL"`
	DEPBaseURL         string `env:"NANODEP_URL"`
	MDMUsername        string `env:"NANOMDM_USERNAME"`
	MDMPassword        string `env:"NANOMDM_PASSWORD"`
	DEPUsername        string `env:"NANODEP_USERNAME"`
	DEPPassword        string `env:"NANODEP_PASSWORD"`
	SudoPassword       string `env:"SUDO_PASSWORD"`
	DEPSyncerContainer string `env:"NANODEP_SYNCER_CONTAINER_NAME"`
	DEPServerName      string `env:"DEP_SERVER_NAME" env-default:"mdm-dep-server"`
}

// SeedConfig holds default user seeding configuration
type SeedConfig struct {
	AdminUsername string `env:"SEED_ADMIN_USERNAME"`
	AdminPassword string `env:"SEED_ADMIN_PASSWORD"`
	AdminEmail    string `env:"SEED_ADMIN_EMAIL"`
}

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
	Cookie   CookieConfig
	Casbin   CasbinConfig
	NanoCMD  NanoCMDConfig
	NanoMDM  NanoMDMConfig
	Seed     SeedConfig
	Redis    RedisConfig

	CORSAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS"`
}

var AppConfig *Config

// Load loads all configuration from environment variables
func Load() (*Config, error) {
	var cfg Config

	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		tlog.Warn("No .env file found, using system environment variables")
	}

	// Read environment variables into Config struct
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Config error: %v", err)
		return nil, err
	}

	AppConfig = &cfg
	return AppConfig, nil
}

// Helper methods on Config

func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

func GetConfig() *Config {
	return AppConfig
}
