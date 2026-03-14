package core

import (
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
	v.SetDefault("database.type", "sqlite")
	v.SetDefault("database.path", "./wsinspect.db")
	v.SetDefault("proxy.target", "ws://localhost:3000")
	v.SetDefault("proxy.buffer_size", 8192)
	v.SetDefault("session.retention_days", 30)
	v.SetDefault("log.level", "info")

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
