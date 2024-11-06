package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const AppConfigPath string = "./config.yaml"

type AppConfig struct {
	GameDir string `yaml:"GameDir"`
}

func ReadAppConfig(filePath string) (*AppConfig, error) {
	var config AppConfig
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error unmarhsalling config file: %v", err)
	}
	return &config, nil
}

func WriteAppConfig(filePath string, config *AppConfig) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}
	defer file.Close()

	if err := yaml.NewEncoder(file).Encode(config); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}
	return nil
}
