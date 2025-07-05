package main

import (
	"testing"

	"conv/cmd"
	"conv/internal/config"
	"conv/internal/currency"
)

func TestParseLegacyArgs(t *testing.T) {
	// Store original function to restore after tests
	originalUserConfigDir := config.UserConfigDirFunc
	defer func() {
		config.UserConfigDirFunc = originalUserConfigDir
	}()

	tests := []struct {
		name            string
		args            []string
		defaultCurrency currency.Currency
		want            currency.Input
		wantErr         bool
	}{
		{
			name:    "valid input with explicit target",
			args:    []string{"100", "USD", "EUR"},
			want:    currency.Input{Amount: 100, From: currency.USD, To: currency.EUR},
			wantErr: false,
		},
		{
			name:            "valid input with default currency",
			args:            []string{"100", "USD"},
			defaultCurrency: currency.EUR,
			want:            currency.Input{Amount: 100, From: currency.USD, To: currency.EUR},
			wantErr:         false,
		},
		{
			name:    "two args but no default currency",
			args:    []string{"100", "USD"},
			want:    currency.Input{},
			wantErr: true,
		},
		{
			name:    "invalid number of arguments - too few",
			args:    []string{"100"},
			want:    currency.Input{},
			wantErr: true,
		},
		{
			name:    "invalid number of arguments - too many",
			args:    []string{"100", "USD", "EUR", "extra"},
			want:    currency.Input{},
			wantErr: true,
		},
		{
			name:    "invalid amount - not a number",
			args:    []string{"not-a-number", "USD", "EUR"},
			want:    currency.Input{},
			wantErr: true,
		},
		{
			name:    "decimal amount with explicit target",
			args:    []string{"100.50", "USD", "EUR"},
			want:    currency.Input{Amount: 100.50, From: currency.USD, To: currency.EUR},
			wantErr: false,
		},
		{
			name:            "decimal amount with default currency",
			args:            []string{"100.50", "USD"},
			defaultCurrency: currency.BRL,
			want:            currency.Input{Amount: 100.50, From: currency.USD, To: currency.BRL},
			wantErr:         false,
		},
		{
			name:    "unsupported source currency",
			args:    []string{"100", "XYZ", "EUR"},
			want:    currency.Input{},
			wantErr: true,
		},
		{
			name:    "unsupported target currency",
			args:    []string{"100", "USD", "XYZ"},
			want:    currency.Input{},
			wantErr: true,
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
			if tt.defaultCurrency != "" {
				err := config.SetDefaultCurrency(string(tt.defaultCurrency))
				if err != nil {
					t.Fatalf("failed to set default currency: %v", err)
				}
			}

			got, err := cmd.ParseLegacyArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLegacyArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Amount != tt.want.Amount {
					t.Errorf("ParseLegacyArgs() Amount = %v, want %v", got.Amount, tt.want.Amount)
				}
				if got.From != tt.want.From {
					t.Errorf("ParseLegacyArgs() From = %v, want %v", got.From, tt.want.From)
				}
				if got.To != tt.want.To {
					t.Errorf("ParseLegacyArgs() To = %v, want %v", got.To, tt.want.To)
				}
			}
		})
	}
}
