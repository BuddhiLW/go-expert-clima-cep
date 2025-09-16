package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCEP(t *testing.T) {
	tests := []struct {
		name     string
		cep      string
		expected bool
	}{
		{
			name:     "CEP válido com 8 dígitos",
			cep:      "01310-100",
			expected: true,
		},
		{
			name:     "CEP válido sem hífen",
			cep:      "01310100",
			expected: true,
		},
		{
			name:     "CEP inválido com menos de 8 dígitos",
			cep:      "1234567",
			expected: false,
		},
		{
			name:     "CEP inválido com mais de 8 dígitos",
			cep:      "123456789",
			expected: false,
		},
		{
			name:     "CEP inválido com letras",
			cep:      "1234567a",
			expected: false,
		},
		{
			name:     "CEP inválido vazio",
			cep:      "",
			expected: false,
		},
		{
			name:     "CEP inválido com caracteres especiais",
			cep:      "123-456-7",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateCEP(tt.cep)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatCEP(t *testing.T) {
	tests := []struct {
		name     string
		cep      string
		expected string
	}{
		{
			name:     "CEP com hífen",
			cep:      "01310-100",
			expected: "01310100",
		},
		{
			name:     "CEP sem hífen",
			cep:      "01310100",
			expected: "01310100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatCEP(tt.cep)
			assert.Equal(t, tt.expected, result)
		})
	}
}
