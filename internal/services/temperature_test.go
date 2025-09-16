package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{
			name:     "0°C para Fahrenheit",
			celsius:  0,
			expected: 32,
		},
		{
			name:     "28.5°C para Fahrenheit",
			celsius:  28.5,
			expected: 83.3,
		},
		{
			name:     "100°C para Fahrenheit",
			celsius:  100,
			expected: 212,
		},
		{
			name:     "-40°C para Fahrenheit",
			celsius:  -40,
			expected: -40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertCelsiusToFahrenheit(tt.celsius)
			assert.InDelta(t, tt.expected, result, 0.1)
		})
	}
}

func TestConvertCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{
			name:     "0°C para Kelvin",
			celsius:  0,
			expected: 273,
		},
		{
			name:     "28.5°C para Kelvin",
			celsius:  28.5,
			expected: 301.5,
		},
		{
			name:     "100°C para Kelvin",
			celsius:  100,
			expected: 373,
		},
		{
			name:     "-273°C para Kelvin",
			celsius:  -273,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertCelsiusToKelvin(tt.celsius)
			assert.InDelta(t, tt.expected, result, 0.1)
		})
	}
}

func TestConvertTemperatures(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected struct {
			fahrenheit float64
			kelvin     float64
		}
	}{
		{
			name:    "28.5°C",
			celsius: 28.5,
			expected: struct {
				fahrenheit float64
				kelvin     float64
			}{
				fahrenheit: 83.3,
				kelvin:     301.5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fahrenheit, kelvin := convertTemperatures(tt.celsius)
			assert.InDelta(t, tt.expected.fahrenheit, fahrenheit, 0.1)
			assert.InDelta(t, tt.expected.kelvin, kelvin, 0.1)
		})
	}
}
