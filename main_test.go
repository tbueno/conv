package main

import (
	"testing"

	"conv/cmd"
	"conv/internal/currency"
)

func TestParseLegacyArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    currency.Input
		wantErr bool
	}{
		{
			name:    "valid input",
			args:    []string{"100", "USD", "EUR"},
			want:    currency.Input{Amount: 100, From: currency.USD, To: currency.EUR},
			wantErr: false,
		},
		{
			name:    "invalid number of arguments - too few",
			args:    []string{"100", "USD"},
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
			name:    "decimal amount",
			args:    []string{"100.50", "USD", "EUR"},
			want:    currency.Input{Amount: 100.50, From: currency.USD, To: currency.EUR},
			wantErr: false,
		},
		{
			name:    "unsupported currency",
			args:    []string{"100", "XYZ", "EUR"},
			want:    currency.Input{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
