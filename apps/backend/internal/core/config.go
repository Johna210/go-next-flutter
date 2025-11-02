package core

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/johna210/go-next-flutter/pkg/utils"
)

func NewConfig() (*Config, error) {
	v := viper.New()

	// Set default config path
	v.SetConfigFile("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	// Read the main config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Read environment specific config file
	env := v.GetString("app.environment")
	if env != "" {
		v.SetConfigName(fmt.Sprintf("config.%s", env))
		// It's okay if environment-specific config doesn't exist
		if err := v.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("error reading %s config file: %w", env, err)
		}
	}

	// Read from environment variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Unmarshal into Config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}
	if c.App.Port <= 0 || c.App.Port > 65535 {
		return fmt.Errorf("app.port must be between 1 and 65535")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("database.dbname is required")
	}

	// Validate App.Environment (used directly in exec.Command)
	if err := utils.IsSafeString(c.App.Environment); err != nil {
		return fmt.Errorf("invalid app.environment: %w", err)
	}

	// Validate DSN components (used in GetDSN, which is then used in exec.Command)
	if err := utils.IsSafeDSNComponent(c.Database.Host); err != nil {
		return fmt.Errorf("invalid database.host: %w", err)
	}
	if err := utils.IsSafeString(c.Database.DBName); err != nil {
		return fmt.Errorf("invalid database.dbname: %w", err)
	}
	if err := utils.IsSafeString(c.Database.SSLMode); err != nil {
		return fmt.Errorf("invalid database.sslmode: %w", err)
	}

	return nil
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Cache.Host, c.Cache.Port)
}
