package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	// Config holds all configuration
	Config Configuration
)

type Configuration struct {
	Kafka KafkaConfig
}

type KafkaConfig struct {
	MessageMaxBytes int `mapstructure:"message_max_bytes"`
}

// GetProjectRoot returns the absolute path to the project root directory
func GetProjectRoot() string {
	// By default, look for the 'backyard' directory as project root
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %s", err)
		return ""
	}

	// Walk up until we find the project root (where go.mod exists)
	dir := currentDir
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the root without finding go.mod
			log.Printf("Could not find project root, using current directory")
			return currentDir
		}
		dir = parent
	}
}

func init() {
	// Set defaults
	viper.SetDefault("kafka.message_max_bytes", 10485760) // 10MB default

	projectRoot := GetProjectRoot()
	configPath := filepath.Join(projectRoot, "config")

	// Look for config file
	viper.SetConfigName("settings") // Name of config file (without extension)
	viper.SetConfigType("yaml")     // YAML format
	viper.AddConfigPath(configPath) // Look in the config directory

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("No config file found, using defaults")
		} else {
			log.Printf("Error reading config file: %s", err)
		}
	}

	// Read config into struct
	if err := viper.Unmarshal(&Config); err != nil {
		log.Printf("Error unmarshaling config: %s", err)
	}
}
