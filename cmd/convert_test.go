package cmd

import (
	"testing"

	"conv/internal/currency"
)

func TestParseConvertArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want currency.Input
	}{
		{
			name: "valid input",
			args: []string{"100", "USD", "EUR"},
			want: currency.Input{Amount: 100, From: currency.USD, To: currency.EUR},
		},
		{
			name: "decimal amount",
			args: []string{"100.50", "USD", "EUR"},
			want: currency.Input{Amount: 100.50, From: currency.USD, To: currency.EUR},
		},
		{
			name: "lowercase currencies",
			args: []string{"50", "usd", "eur"},
			want: currency.Input{Amount: 50, From: currency.USD, To: currency.EUR},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseConvertArgs(tt.args)
			if got.Amount != tt.want.Amount {
				t.Errorf("parseConvertArgs() Amount = %v, want %v", got.Amount, tt.want.Amount)
			}
			if got.From != tt.want.From {
				t.Errorf("parseConvertArgs() From = %v, want %v", got.From, tt.want.From)
			}
			if got.To != tt.want.To {
				t.Errorf("parseConvertArgs() To = %v, want %v", got.To, tt.want.To)
			}
		})
	}
}

func TestValidateConvertArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid arguments",
			args:    []string{"100", "USD", "EUR"},
			wantErr: false,
		},
		{
			name:    "invalid amount - not a number",
			args:    []string{"not-a-number", "USD", "EUR"},
			wantErr: true,
		},
		{
			name:    "unsupported source currency",
			args:    []string{"100", "INVALID", "EUR"},
			wantErr: true,
		},
		{
			name:    "unsupported target currency",
			args:    []string{"100", "USD", "INVALID"},
			wantErr: true,
		},
		{
			name:    "too few arguments",
			args:    []string{"100", "USD"},
			wantErr: true,
		},
		{
			name:    "too many arguments",
			args:    []string{"100", "USD", "EUR", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConvertArgs(nil, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConvertArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}