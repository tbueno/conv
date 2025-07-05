package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"conv/internal/config"
)

func TestConfigCommand(t *testing.T) {
	// Store original function to restore after tests
	originalUserConfigDir := config.UserConfigDirFunc
	defer func() {
		config.UserConfigDirFunc = originalUserConfigDir
	}()

	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		wantOutputContains []string
	}{
		{
			name:    "config without subcommand shows error",
			args:    []string{},
			wantErr: true,
			wantOutputContains: []string{"requires at least 1 arg(s)"},
		},
		{
			name:    "config set without arguments shows help",
			args:    []string{"set"},
			wantErr: true,
		},
		{
			name:    "config get without arguments shows help", 
			args:    []string{"get"},
			wantErr: true,
		},
		{
			name:    "config show displays empty config",
			args:    []string{"show"},
			wantErr: false,
			wantOutputContains: []string{"Configuration:", "Default currency: (not set)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh temp directory for each test
			testTempDir := t.TempDir()
			
			// Reset global config and set new temp dir for each test
			config.ResetGlobalConfig()
			config.UserConfigDirFunc = func() (string, error) {
				return testTempDir, nil
			}

			// Capture output
			var buf bytes.Buffer
			
			// Create a new command instance for testing
			cmd := &cobra.Command{Use: "test"}
			cmd.AddCommand(configCmd)
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			
			// Set arguments
			fullArgs := append([]string{"config"}, tt.args...)
			cmd.SetArgs(fullArgs)
			
			// Execute command
			err := cmd.Execute()
			
			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("config command error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			// Check output contains expected strings
			output := buf.String()
			for _, expectedOutput := range tt.wantOutputContains {
				if !strings.Contains(output, expectedOutput) {
					t.Errorf("expected output to contain %q, got: %s", expectedOutput, output)
				}
			}
		})
	}
}

func TestConfigSetCommand(t *testing.T) {
	// Store original function to restore after tests
	originalUserConfigDir := config.UserConfigDirFunc
	defer func() {
		config.UserConfigDirFunc = originalUserConfigDir
	}()

	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		wantOutputContains []string
	}{
		{
			name:    "set valid default currency",
			args:    []string{"set", "default-currency", "USD"},
			wantErr: false,
			wantOutputContains: []string{"Default currency set to: USD"},
		},
		{
			name:    "set invalid default currency",
			args:    []string{"set", "default-currency", "INVALID"},
			wantErr: false,
			wantOutputContains: []string{"Error setting default currency", "unsupported currency"},
		},
		{
			name:    "set unknown setting",
			args:    []string{"set", "unknown-setting", "value"},
			wantErr: false,
			wantOutputContains: []string{"Error: unknown setting", "Available settings: default-currency"},
		},
		{
			name:    "clear default currency",
			args:    []string{"set", "default-currency", "clear"},
			wantErr: false,
			wantOutputContains: []string{"Default currency cleared"},
		},
		{
			name:    "clear default currency with none",
			args:    []string{"set", "default-currency", "none"},
			wantErr: false,
			wantOutputContains: []string{"Default currency cleared"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh temp directory for each test
			testTempDir := t.TempDir()
			
			// Reset global config and set new temp dir for each test
			config.ResetGlobalConfig()
			config.UserConfigDirFunc = func() (string, error) {
				return testTempDir, nil
			}

			// Capture output
			var buf bytes.Buffer
			
			// Create a new command instance for testing
			cmd := &cobra.Command{Use: "test"}
			cmd.AddCommand(configCmd)
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			
			// Set arguments
			fullArgs := append([]string{"config"}, tt.args...)
			cmd.SetArgs(fullArgs)
			
			// Execute command
			err := cmd.Execute()
			
			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("config set command error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			// Check output contains expected strings
			output := buf.String()
			for _, expectedOutput := range tt.wantOutputContains {
				if !strings.Contains(output, expectedOutput) {
					t.Errorf("expected output to contain %q, got: %s", expectedOutput, output)
				}
			}
		})
	}
}

func TestConfigGetCommand(t *testing.T) {
	// Store original function to restore after tests
	originalUserConfigDir := config.UserConfigDirFunc
	defer func() {
		config.UserConfigDirFunc = originalUserConfigDir
	}()

	tests := []struct {
		name               string
		args               []string
		setupDefaultCurrency string
		wantErr            bool
		wantOutputContains []string
	}{
		{
			name:               "get default currency when set",
			args:               []string{"get", "default-currency"},
			setupDefaultCurrency: "EUR",
			wantErr:            false,
			wantOutputContains: []string{"Default currency: EUR"},
		},
		{
			name:    "get default currency when not set",
			args:    []string{"get", "default-currency"},
			wantErr: false,
			wantOutputContains: []string{"No default currency set"},
		},
		{
			name:    "get unknown setting",
			args:    []string{"get", "unknown-setting"},
			wantErr: false,
			wantOutputContains: []string{"Error: unknown setting", "Available settings: default-currency"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh temp directory for each test
			testTempDir := t.TempDir()
			
			// Reset global config and set new temp dir for each test
			config.ResetGlobalConfig()
			config.UserConfigDirFunc = func() (string, error) {
				return testTempDir, nil
			}

			// Set up default currency if needed
			if tt.setupDefaultCurrency != "" {
				err := config.SetDefaultCurrency(tt.setupDefaultCurrency)
				if err != nil {
					t.Fatalf("failed to set up default currency: %v", err)
				}
			}

			// Capture output
			var buf bytes.Buffer
			
			// Create a new command instance for testing
			cmd := &cobra.Command{Use: "test"}
			cmd.AddCommand(configCmd)
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			
			// Set arguments
			fullArgs := append([]string{"config"}, tt.args...)
			cmd.SetArgs(fullArgs)
			
			// Execute command
			err := cmd.Execute()
			
			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("config get command error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			// Check output contains expected strings
			output := buf.String()
			for _, expectedOutput := range tt.wantOutputContains {
				if !strings.Contains(output, expectedOutput) {
					t.Errorf("expected output to contain %q, got: %s", expectedOutput, output)
				}
			}
		})
	}
}

func TestConfigShowCommand(t *testing.T) {
	// Store original function to restore after tests
	originalUserConfigDir := config.UserConfigDirFunc
	defer func() {
		config.UserConfigDirFunc = originalUserConfigDir
	}()

	tests := []struct {
		name               string
		setupDefaultCurrency string
		wantOutputContains []string
	}{
		{
			name: "show config when not set",
			wantOutputContains: []string{"Configuration:", "Default currency: (not set)"},
		},
		{
			name:               "show config when set to USD",
			setupDefaultCurrency: "USD",
			wantOutputContains: []string{"Configuration:", "Default currency: USD"},
		},
		{
			name:               "show config when set to EUR",
			setupDefaultCurrency: "EUR",
			wantOutputContains: []string{"Configuration:", "Default currency: EUR"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh temp directory for each test
			testTempDir := t.TempDir()
			
			// Reset global config and set new temp dir for each test
			config.ResetGlobalConfig()
			config.UserConfigDirFunc = func() (string, error) {
				return testTempDir, nil
			}

			// Set up default currency if needed
			if tt.setupDefaultCurrency != "" {
				err := config.SetDefaultCurrency(tt.setupDefaultCurrency)
				if err != nil {
					t.Fatalf("failed to set up default currency: %v", err)
				}
			}

			// Capture output
			var buf bytes.Buffer
			
			// Create a new command instance for testing
			cmd := &cobra.Command{Use: "test"}
			cmd.AddCommand(configCmd)
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			
			// Set arguments
			cmd.SetArgs([]string{"config", "show"})
			
			// Execute command
			err := cmd.Execute()
			if err != nil {
				t.Errorf("config show command unexpected error: %v", err)
				return
			}
			
			// Check output contains expected strings
			output := buf.String()
			for _, expectedOutput := range tt.wantOutputContains {
				if !strings.Contains(output, expectedOutput) {
					t.Errorf("expected output to contain %q, got: %s", expectedOutput, output)
				}
			}
		})
	}
}