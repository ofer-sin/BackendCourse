package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	// Read the configuration file
	// viper.SetConfigFile(path)
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // Specify the config file type, e.g., "json", "yaml", "toml", etc.

	// read environment variables that override config file
	viper.AutomaticEnv()

	fmt.Println("Connecting to database with source", config.DBSource)

	// Read the configuration
	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	// Unmarshal the configuration into the Config struct
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	fmt.Println("After Unmarshal")
	fmt.Println("Starting server on port", config.ServerAddress)
	fmt.Println("Connecting to database with driver", config.DBDriver)
	fmt.Println("Connecting to database with source", config.DBSource)

	return config, nil
}
