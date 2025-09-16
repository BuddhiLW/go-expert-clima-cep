package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"cep-temperatura/internal/models"
)

// WeatherService interface para operações de clima
type WeatherService interface {
	GetTemperature(city, state string) (float64, error)
}

type weatherService struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewWeatherService cria uma nova instância do serviço de clima
func NewWeatherService() WeatherService {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		apiKey = "b5d4215a52bf4e2da2f144209251609" // Chave padrão
	}

	return &weatherService{
		baseURL: "http://api.weatherapi.com/v1",
		apiKey:  apiKey,
		client:  &http.Client{},
	}
}

// GetTemperature busca a temperatura atual de uma cidade
func (s *weatherService) GetTemperature(city, state string) (float64, error) {
	// Construir query para a API
	query := fmt.Sprintf("%s, %s, Brazil", city, state)
	encodedQuery := url.QueryEscape(query)
	apiURL := fmt.Sprintf("%s/current.json?key=%s&q=%s", s.baseURL, s.apiKey, encodedQuery)

	resp, err := s.client.Get(apiURL)
	if err != nil {
		return 0, fmt.Errorf("erro ao consultar clima: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("erro ao consultar clima: status %d", resp.StatusCode)
	}

	var weatherResponse models.WeatherResponse
	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		return 0, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return weatherResponse.Current.TempC, nil
}
