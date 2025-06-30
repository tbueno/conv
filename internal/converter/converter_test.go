package converter

import (
	"errors"
	"testing"

	"conv/internal/currency"
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
		from          currency.Currency
		to            currency.Currency
		mockConverter *MockConverter
		want          float32
		wantErr       bool
	}{
		{
			name:          "successful conversion",
			amount:        100,
			from:          currency.USD,
			to:            currency.EUR,
			mockConverter: &MockConverter{shouldError: false, value: 85.5},
			want:          85.5,
			wantErr:       false,
		},
		{
			name:          "converter error",
			amount:        100,
			from:          currency.USD,
			to:            currency.EUR,
			mockConverter: &MockConverter{shouldError: true},
			want:          0,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := currency.Input{Amount: tt.amount, From: tt.from, To: tt.to}
			got, err := Convert(input, tt.mockConverter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}