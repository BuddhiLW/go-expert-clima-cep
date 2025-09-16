package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cep-temperatura/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCEPService é um mock do CEPService
type MockCEPService struct {
	mock.Mock
}

func (m *MockCEPService) ValidateCEP(cep string) bool {
	args := m.Called(cep)
	return args.Bool(0)
}

func (m *MockCEPService) GetLocation(cep string) (*models.CEPResponse, error) {
	args := m.Called(cep)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CEPResponse), args.Error(1)
}

// MockWeatherService é um mock do WeatherService
type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetTemperature(city, state string) (float64, error) {
	args := m.Called(city, state)
	return args.Get(0).(float64), args.Error(1)
}

// MockTemperatureService é um mock do TemperatureService
type MockTemperatureService struct {
	mock.Mock
}

func (m *MockTemperatureService) ConvertTemperatures(celsius float64) (fahrenheit, kelvin float64) {
	args := m.Called(celsius)
	return args.Get(0).(float64), args.Get(1).(float64)
}

func TestTemperatureHandler_GetTemperature_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mocks
	mockCEPService := new(MockCEPService)
	mockWeatherService := new(MockWeatherService)
	mockTemperatureService := new(MockTemperatureService)

	// Configurar mocks
	mockCEPService.On("ValidateCEP", "01310100").Return(true)
	mockCEPService.On("GetLocation", "01310100").Return(&models.CEPResponse{
		Localidade: "São Paulo",
		UF:         "SP",
	}, nil)

	mockWeatherService.On("GetTemperature", "São Paulo", "SP").Return(28.5, nil)
	mockTemperatureService.On("ConvertTemperatures", 28.5).Return(83.3, 301.5)

	// Criar handler
	handler := NewTemperatureHandler(mockCEPService, mockWeatherService, mockTemperatureService)

	// Criar request
	req, _ := http.NewRequest("GET", "/temperature/01310100", nil)
	w := httptest.NewRecorder()

	// Criar contexto Gin
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "cep", Value: "01310100"}}

	// Executar handler
	handler.GetTemperature(c)

	// Verificar resposta
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.TemperatureResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 28.5, response.TempC)
	assert.Equal(t, 83.3, response.TempF)
	assert.Equal(t, 301.5, response.TempK)

	// Verificar se todos os mocks foram chamados
	mockCEPService.AssertExpectations(t)
	mockWeatherService.AssertExpectations(t)
	mockTemperatureService.AssertExpectations(t)
}

func TestTemperatureHandler_GetTemperature_InvalidCEP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mocks
	mockCEPService := new(MockCEPService)
	mockWeatherService := new(MockWeatherService)
	mockTemperatureService := new(MockTemperatureService)

	// Configurar mock para CEP inválido
	mockCEPService.On("ValidateCEP", "123").Return(false)

	// Criar handler
	handler := NewTemperatureHandler(mockCEPService, mockWeatherService, mockTemperatureService)

	// Criar request
	req, _ := http.NewRequest("GET", "/temperature/123", nil)
	w := httptest.NewRecorder()

	// Criar contexto Gin
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "cep", Value: "123"}}

	// Executar handler
	handler.GetTemperature(c)

	// Verificar resposta
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Contains(t, w.Body.String(), "invalid zipcode")

	// Verificar se apenas o mock de validação foi chamado
	mockCEPService.AssertExpectations(t)
}

func TestTemperatureHandler_GetTemperature_CEPNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mocks
	mockCEPService := new(MockCEPService)
	mockWeatherService := new(MockWeatherService)
	mockTemperatureService := new(MockTemperatureService)

	// Configurar mocks
	mockCEPService.On("ValidateCEP", "99999999").Return(true)
	mockCEPService.On("GetLocation", "99999999").Return(nil, assert.AnError)

	// Criar handler
	handler := NewTemperatureHandler(mockCEPService, mockWeatherService, mockTemperatureService)

	// Criar request
	req, _ := http.NewRequest("GET", "/temperature/99999999", nil)
	w := httptest.NewRecorder()

	// Criar contexto Gin
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "cep", Value: "99999999"}}

	// Executar handler
	handler.GetTemperature(c)

	// Verificar resposta
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "can not find zipcode")

	// Verificar se os mocks foram chamados
	mockCEPService.AssertExpectations(t)
}
