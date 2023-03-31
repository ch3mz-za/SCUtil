package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const AppConfigPath string = "./config.toml"

type AppConfig struct {
	GameDir string
}

func ReadAppConfig(filePath string) (*AppConfig, error) {
	var config AppConfig
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		return nil, fmt.Errorf("error loading config file: %v", err)
	}
	return &config, nil
}

func WriteAppConfig(filePath string, config *AppConfig) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(config); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}
	return nil
}
