package kjspkg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func HasConfig(path string) bool {
	configFilePath := filepath.Join(path, ".kjspkg")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetConfig(path string, noKubeOk bool) (*Config, error) {
	config := DefaultConfig()
	configFilePath := filepath.Join(path, ".kjspkg")

	if !noKubeOk && !IsKube(path) {
		return nil, fmt.Errorf("are you sure this is the kubejs directory?")
	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf(".kjspkg file not found")
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read .kjspkg file: %w", err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse .kjspkg file as JSON: %w", err)
	}

	return config, nil
}

func SetConfig(path string, config *Config) error {
	configFilePath := filepath.Join(path, ".kjspkg")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	if err := os.WriteFile(configFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write .kjspkg: %w", err)
	}

	return nil
}
