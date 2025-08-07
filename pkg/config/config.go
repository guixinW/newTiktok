package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

// Config holds all configuration for the application.
type Config struct {
	Port        string      `mapstructure:"port"`
	LogLevel    string      `mapstructure:"log_level"`
	Database    Database    `mapstructure:"database"`
	Redis       Redis       `mapstructure:"redis"`
	Gateway     Gateway     `mapstructure:"gateway"`
	UserService UserService `mapstructure:"user_service"`
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

// Gateway holds all gateway configuration.
type Gateway struct {
	Port string `mapstructure:"port"`
}

// UserService holds user service configuration.
type UserService struct {
	Port        string `mapstructure:"port"`
	ServiceName string `mapstructure:"service_name"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	// For containerized environments, we specify the full path.
	// The Kubernetes manifest mounts the config file to /app/config.yaml
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml" // 默认在当前目录下查找 config.yaml
	}
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// ReadInConfig will now use the path from SetConfigFile if it exists,
	// otherwise it will search in the paths from AddConfigPath.
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
