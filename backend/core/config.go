package core

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	v *viper.Viper
}

func NewConfig() *Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("$HOME/.wsinspect")

	// Set default values
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("database.type", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "wsinspect")
	v.SetDefault("database.password", "wsinspect_password")
	v.SetDefault("database.name", "wsinspect_db")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("proxy.target", "ws://localhost:3000")
	v.SetDefault("proxy.buffer_size", 8192)
	v.SetDefault("session.retention_days", 30)
	v.SetDefault("log.level", "info")

	// Environment variables support
	v.BindEnv("database.url", "DATABASE_URL")
	v.BindEnv("server.port", "SERVER_PORT", "PORT")
	v.BindEnv("server.host", "SERVER_HOST")
	v.BindEnv("database.host", "DB_HOST", "POSTGRES_HOST")
	v.BindEnv("database.port", "DB_PORT", "POSTGRES_PORT")
	v.BindEnv("database.user", "DB_USER", "POSTGRES_USER")
	v.BindEnv("database.password", "DB_PASSWORD", "POSTGRES_PASSWORD")
	v.BindEnv("database.name", "DB_NAME", "POSTGRES_DB")
	v.BindEnv("database.sslmode", "DB_SSLMODE", "POSTGRES_SSLMODE")
	v.BindEnv("proxy.target", "PROXY_TARGET")
	v.BindEnv("log.level", "LOG_LEVEL")

	// Try to read config file (optional)
	if err := v.ReadInConfig(); err != nil {
		log.Printf("No config file found, using defaults and environment variables: %v", err)
	}

	return &Config{v: v}
}

func (c *Config) Load() error {
	return c.v.ReadInConfig()
}

func (c *Config) GetString(key string) string {
	return c.v.GetString(key)
}

func (c *Config) GetInt(key string) int {
	return c.v.GetInt(key)
}

func (c *Config) GetBool(key string) bool {
	return c.v.GetBool(key)
}

func (c *Config) GetStringSlice(key string) []string {
	return c.v.GetStringSlice(key)
}

// GetDatabaseURL returns the full database URL from environment or builds it from individual settings
func (c *Config) GetDatabaseURL() string {
	// Check if DATABASE_URL is set directly
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	// Build from individual settings
	host := c.GetString("database.host")
	port := c.GetInt("database.port")
	user := c.GetString("database.user")
	password := c.GetString("database.password")
	dbname := c.GetString("database.name")
	sslmode := c.GetString("database.sslmode")

	return "postgres://" + user + ":" + password + "@" + host + ":" + strconv.Itoa(port) + "/" + dbname + "?sslmode=" + sslmode
}


