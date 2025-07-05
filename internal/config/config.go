package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"conv/internal/currency"
)

type Config struct {
	DefaultCurrency currency.Currency `json:"default_currency,omitempty"`
}

var globalConfig *Config

// userConfigDirFunc allows mocking os.UserConfigDir in tests
var userConfigDirFunc = os.UserConfigDir

func getConfigFilePath() (string, error) {
	configDir, err := userConfigDirFunc()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}
	
	appConfigDir := filepath.Join(configDir, "conv")
	err = os.MkdirAll(appConfigDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}
	
	return filepath.Join(appConfigDir, "config.json"), nil
}

func LoadConfig() (*Config, error) {
	if globalConfig != nil {
		return globalConfig, nil
	}
	
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	
	config := &Config{}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			globalConfig = config
			return config, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	globalConfig = config
	return config, nil
}

func SaveConfig(config *Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	globalConfig = config
	return nil
}

func SetDefaultCurrency(currencyCode string) error {
	curr := currency.Currency(strings.ToUpper(currencyCode))
	if !curr.IsValid() {
		return fmt.Errorf("unsupported currency: %s", currencyCode)
	}
	
	config, err := LoadConfig()
	if err != nil {
		return err
	}
	
	config.DefaultCurrency = curr
	return SaveConfig(config)
}

func GetDefaultCurrency() (currency.Currency, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}
	
	return config.DefaultCurrency, nil
}

func GetConfig() (*Config, error) {
	return LoadConfig()
}