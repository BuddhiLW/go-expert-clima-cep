package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cep-temperatura/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestWeatherService_GetTemperature(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simular resposta da WeatherAPI
		response := models.WeatherResponse{
			Location: struct {
				Name    string `json:"name"`
				Region  string `json:"region"`
				Country string `json:"country"`
			}{
				Name:    "São Paulo",
				Region:  "São Paulo",
				Country: "Brazil",
			},
			Current: struct {
				TempC float64 `json:"temp_c"`
			}{
				TempC: 28.5,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Criar serviço com URL do mock
	weatherService := &weatherService{
		baseURL: server.URL,
		apiKey:  "b5d4215a52bf4e2da2f144209251609",
		client:  &http.Client{},
	}

	t.Run("busca temperatura com sucesso", func(t *testing.T) {
		temp, err := weatherService.GetTemperature("São Paulo", "SP")
		if err != nil {
			t.Logf("Erro: %v", err)
		}
		assert.NoError(t, err)
		assert.Equal(t, 28.5, temp)
	})
}

func TestWeatherService_GetTemperature_Error(t *testing.T) {
	// Mock server que retorna erro
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": {"message": "No matching location found."}}`))
	}))
	defer server.Close()

	weatherService := &weatherService{
		baseURL: server.URL,
		apiKey:  "b5d4215a52bf4e2da2f144209251609",
		client:  &http.Client{},
	}

	t.Run("erro ao buscar temperatura", func(t *testing.T) {
		_, err := weatherService.GetTemperature("CidadeInexistente", "XX")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao consultar clima")
	})
}
