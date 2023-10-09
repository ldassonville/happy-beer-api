package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func InitConfigs(configs any) any {
	path := os.Getenv("CONFIGURATION_PATH")

	// Use conventional config path
	if path == "" {
		path = "/workspace/"
	}

	return InitConfFromPath(path, "config", configs)
}

func InitSecrets(secrets any) any {
	path := os.Getenv("SECRET_PATH")

	// Use conventional config path
	if path == "" {
		path = "/workspace/"
	}

	return InitConfFromPath(path, "secret", secrets)
}

func InitConfFromPath(path string, name string, config any) any {

	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		logrus.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}
	err = viper.Unmarshal(config)
	if err != nil { // Handle errors reading the config file
		logrus.Fatal(fmt.Errorf("fatal error while unmarshaling config file: %w", err))
	}
	return config
}
