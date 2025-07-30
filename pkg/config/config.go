package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config holds all configuration for the application.
type Config struct {
	Port     string   `mapstructure:"port"`
	LogLevel string   `mapstructure:"log_level"`
	Database Database `mapstructure:"database"`
	Redis    Redis    `mapstructure:"redis"`
}

// Database holds all database configuration.
type Database struct {
	DSN string `mapstructure:"dsn"`
}

// Redis holds all Redis configuration.
type Redis struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return
}
