package utils

import (
	"go-parallel_queue/internal/config"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	defaultAppConfigPath = "config/local.yaml"
	envVarAppConfigPath  = "APP_CONFIG_PATH"
)

// ENV -> default
func MustLoadConfig() *config.Config {
	configPath := defaultAppConfigPath

	if fromEnv := os.Getenv(envVarAppConfigPath); fromEnv != "" {
		configPath = fromEnv
	}

	dir, _ := os.Getwd()
	log.Printf("Using config_path: %s, wd: %s", configPath, dir)
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("wrong config (%s): %v", configPath, err)
	}
	return config
}

func loadConfig(configPath string) (*config.Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config config.Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
