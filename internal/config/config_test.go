package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"conv/internal/currency"
)

func TestConfig_LoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		setupConfig    *Config
		wantErr        bool
		wantCurrency   currency.Currency
	}{
		{
			name:           "load empty config when file doesn't exist",
			setupConfig:    nil,
			wantErr:        false,
			wantCurrency:   "",
		},
		{
			name: "load existing config with default currency",
			setupConfig: &Config{
				DefaultCurrency: currency.USD,
			},
			wantErr:        false,
			wantCurrency:   currency.USD,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global config for each test
			globalConfig = nil
			
			// Create temporary config directory
			tempDir := t.TempDir()
			originalUserConfigDir := userConfigDirFunc
			defer func() {
				userConfigDirFunc = originalUserConfigDir
			}()
			
			// Mock userConfigDirFunc to return our temp directory
			userConfigDirFunc = func() (string, error) {
				return tempDir, nil
			}

			// Setup config file if needed
			if tt.setupConfig != nil {
				configPath := filepath.Join(tempDir, "conv", "config.json")
				err := os.MkdirAll(filepath.Dir(configPath), 0755)
				if err != nil {
					t.Fatalf("failed to create config directory: %v", err)
				}

				data, err := json.MarshalIndent(tt.setupConfig, "", "  ")
				if err != nil {
					t.Fatalf("failed to marshal config: %v", err)
				}

				err = os.WriteFile(configPath, data, 0644)
				if err != nil {
					t.Fatalf("failed to write config file: %v", err)
				}
			}

			// Test LoadConfig
			config, err := LoadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && config.DefaultCurrency != tt.wantCurrency {
				t.Errorf("LoadConfig() DefaultCurrency = %v, want %v", config.DefaultCurrency, tt.wantCurrency)
			}
		})
	}
}

func TestConfig_SaveConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "save config with default currency",
			config: &Config{
				DefaultCurrency: currency.EUR,
			},
			wantErr: false,
		},
		{
			name: "save empty config",
			config: &Config{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global config for each test
			globalConfig = nil

			// Create temporary config directory
			tempDir := t.TempDir()
			originalUserConfigDir := userConfigDirFunc
			defer func() {
				userConfigDirFunc = originalUserConfigDir
			}()

			// Mock userConfigDirFunc to return our temp directory
			userConfigDirFunc = func() (string, error) {
				return tempDir, nil
			}

			// Test SaveConfig
			err := SaveConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify config was saved correctly
				configPath := filepath.Join(tempDir, "conv", "config.json")
				data, err := os.ReadFile(configPath)
				if err != nil {
					t.Errorf("failed to read saved config file: %v", err)
					return
				}

				var savedConfig Config
				err = json.Unmarshal(data, &savedConfig)
				if err != nil {
					t.Errorf("failed to parse saved config file: %v", err)
					return
				}

				if savedConfig.DefaultCurrency != tt.config.DefaultCurrency {
					t.Errorf("SaveConfig() saved DefaultCurrency = %v, want %v", savedConfig.DefaultCurrency, tt.config.DefaultCurrency)
				}
			}
		})
	}
}

func TestConfig_SetDefaultCurrency(t *testing.T) {
	tests := []struct {
		name         string
		currencyCode string
		wantErr      bool
		wantCurrency currency.Currency
	}{
		{
			name:         "set valid currency USD",
			currencyCode: "USD",
			wantErr:      false,
			wantCurrency: currency.USD,
		},
		{
			name:         "set valid currency EUR",
			currencyCode: "EUR",
			wantErr:      false,
			wantCurrency: currency.EUR,
		},
		{
			name:         "set valid currency lowercase",
			currencyCode: "usd",
			wantErr:      false,
			wantCurrency: currency.USD,
		},
		{
			name:         "set invalid currency",
			currencyCode: "INVALID",
			wantErr:      true,
			wantCurrency: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global config for each test
			globalConfig = nil

			// Create temporary config directory
			tempDir := t.TempDir()
			originalUserConfigDir := userConfigDirFunc
			defer func() {
				userConfigDirFunc = originalUserConfigDir
			}()

			// Mock userConfigDirFunc to return our temp directory
			userConfigDirFunc = func() (string, error) {
				return tempDir, nil
			}

			// Test SetDefaultCurrency
			err := SetDefaultCurrency(tt.currencyCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDefaultCurrency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify config was set correctly
				config, err := LoadConfig()
				if err != nil {
					t.Errorf("failed to load config after SetDefaultCurrency: %v", err)
					return
				}

				if config.DefaultCurrency != tt.wantCurrency {
					t.Errorf("SetDefaultCurrency() set DefaultCurrency = %v, want %v", config.DefaultCurrency, tt.wantCurrency)
				}
			}
		})
	}
}

func TestConfig_GetDefaultCurrency(t *testing.T) {
	tests := []struct {
		name         string
		setupConfig  *Config
		wantErr      bool
		wantCurrency currency.Currency
	}{
		{
			name: "get default currency when set",
			setupConfig: &Config{
				DefaultCurrency: currency.EUR,
			},
			wantErr:      false,
			wantCurrency: currency.EUR,
		},
		{
			name:         "get default currency when not set",
			setupConfig:  &Config{},
			wantErr:      false,
			wantCurrency: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global config for each test
			globalConfig = nil

			// Create temporary config directory
			tempDir := t.TempDir()
			originalUserConfigDir := userConfigDirFunc
			defer func() {
				userConfigDirFunc = originalUserConfigDir
			}()

			// Mock userConfigDirFunc to return our temp directory
			userConfigDirFunc = func() (string, error) {
				return tempDir, nil
			}

			// Setup config
			if tt.setupConfig != nil {
				err := SaveConfig(tt.setupConfig)
				if err != nil {
					t.Fatalf("failed to setup config: %v", err)
				}
			}

			// Test GetDefaultCurrency
			currency, err := GetDefaultCurrency()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDefaultCurrency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && currency != tt.wantCurrency {
				t.Errorf("GetDefaultCurrency() = %v, want %v", currency, tt.wantCurrency)
			}
		})
	}
}

func TestConfig_GetConfig(t *testing.T) {
	// Reset global config for test
	globalConfig = nil

	// Create temporary config directory
	tempDir := t.TempDir()
	originalUserConfigDir := userConfigDirFunc
	defer func() {
		userConfigDirFunc = originalUserConfigDir
	}()

	// Mock userConfigDirFunc to return our temp directory
	userConfigDirFunc = func() (string, error) {
		return tempDir, nil
	}

	// Test GetConfig
	config, err := GetConfig()
	if err != nil {
		t.Errorf("GetConfig() error = %v", err)
		return
	}

	if config == nil {
		t.Error("GetConfig() returned nil config")
	}
}

func TestConfig_getConfigFilePath(t *testing.T) {
	// Create temporary config directory
	tempDir := t.TempDir()
	originalUserConfigDir := userConfigDirFunc
	defer func() {
		userConfigDirFunc = originalUserConfigDir
	}()

	// Mock userConfigDirFunc to return our temp directory
	userConfigDirFunc = func() (string, error) {
		return tempDir, nil
	}

	configPath, err := getConfigFilePath()
	if err != nil {
		t.Errorf("getConfigFilePath() error = %v", err)
		return
	}

	expectedPath := filepath.Join(tempDir, "conv", "config.json")
	if configPath != expectedPath {
		t.Errorf("getConfigFilePath() = %v, want %v", configPath, expectedPath)
	}

	// Verify directory was created
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("getConfigFilePath() did not create config directory")
	}
}