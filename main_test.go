package main

import (
	"errors"
	"testing"
)

// MockConverter implements the Converter interface for testing
type MockConverter struct {
	shouldError bool
	value       float32
}

func (m *MockConverter) Convert(amount float32, from, to string) (float32, error) {
	if m.shouldError {
		return 0, errors.New("mock error")
	}
	return m.value, nil
}

func TestConvert(t *testing.T) {
	tests := []struct {
		name          string
		amount        float32
		from          Currency
		to            Currency
		mockConverter *MockConverter
		want          float32
		wantErr       bool
	}{
		{
			name:          "successful conversion",
			amount:        100,
			from:          USD,
			to:            EUR,
			mockConverter: &MockConverter{shouldError: false, value: 85.5},
			want:          85.5,
			wantErr:       false,
		},
		{
			name:          "converter error",
			amount:        100,
			from:          USD,
			to:            EUR,
			mockConverter: &MockConverter{shouldError: true},
			want:          0,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convert(Input{Amount: tt.amount, From: tt.from, To: tt.to}, tt.mockConverter)
			if (err != nil) != tt.wantErr {
				t.Errorf("convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    Input
		wantErr bool
	}{
		{
			name:    "valid input",
			args:    []string{"program", "100", "USD", "EUR"},
			want:    Input{Amount: 100, From: USD, To: EUR},
			wantErr: false,
		},
		{
			name:    "invalid number of arguments - too few",
			args:    []string{"program", "100", "USD"},
			want:    Input{},
			wantErr: true,
		},
		{
			name:    "invalid number of arguments - too many",
			args:    []string{"program", "100", "USD", "EUR", "extra"},
			want:    Input{},
			wantErr: true,
		},
		{
			name:    "invalid amount - not a number",
			args:    []string{"program", "not-a-number", "USD", "EUR"},
			want:    Input{},
			wantErr: true,
		},
		{
			name:    "empty arguments",
			args:    []string{"program", "", "", ""},
			want:    Input{},
			wantErr: true,
		},
		{
			name:    "decimal amount",
			args:    []string{"program", "100.50", "USD", "EUR"},
			want:    Input{Amount: 100.50, From: USD, To: EUR},
			wantErr: false,
		},
		{
			name:    "unsupported currency",
			args:    []string{"program", "100", "XYZ", "EUR"},
			want:    Input{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Amount != tt.want.Amount {
					t.Errorf("parseArgs() Amount = %v, want %v", got.Amount, tt.want.Amount)
				}
				if got.From != tt.want.From {
					t.Errorf("parseArgs() From = %v, want %v", got.From, tt.want.From)
				}
				if got.To != tt.want.To {
					t.Errorf("parseArgs() To = %v, want %v", got.To, tt.want.To)
				}
			}
		})
	}
}

func TestCurrencyValidation(t *testing.T) {
	tests := []struct {
		name     string
		currency Currency
		want     bool
	}{
		{
			name:     "hardcoded USD should be valid",
			currency: USD,
			want:     true,
		},
		{
			name:     "hardcoded EUR should be valid",
			currency: EUR,
			want:     true,
		},
		{
			name:     "JPY should be valid (from API cache)",
			currency: "JPY",
			want:     true,
		},
		{
			name:     "invalid currency should be false",
			currency: "INVALID",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.currency.IsValid()
			if got != tt.want {
				t.Errorf("Currency.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
