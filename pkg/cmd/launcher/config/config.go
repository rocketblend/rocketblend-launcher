package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var Name = "settings"
var Type = "json"

func Load(name string) (*viper.Viper, error) {
	v := viper.New()

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("cannot find config directory: %v", err)
	}

	appDir := filepath.Join(configDir, name)

	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create main directory: %w", err)
	}

	v.SetDefault("previous", "")

	v.SetConfigName(Name)   // Set the name of the configuration file
	v.AddConfigPath(appDir) // Look for the configuration file at the home directory
	v.SetConfigType(Type)   // Set the config type to JSON

	v.SafeWriteConfig()

	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return v, nil
}
