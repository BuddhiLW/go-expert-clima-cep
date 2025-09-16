package handlers

import (
	"net/http"

	"cep-temperatura/internal/models"
	"cep-temperatura/internal/services"

	"github.com/gin-gonic/gin"
)

// TemperatureHandler gerencia as requisições de temperatura
type TemperatureHandler struct {
	cepService         services.CEPService
	weatherService     services.WeatherService
	temperatureService services.TemperatureService
}

// NewTemperatureHandler cria uma nova instância do handler de temperatura
func NewTemperatureHandler(
	cepService services.CEPService,
	weatherService services.WeatherService,
	temperatureService services.TemperatureService,
) *TemperatureHandler {
	return &TemperatureHandler{
		cepService:         cepService,
		weatherService:     weatherService,
		temperatureService: temperatureService,
	}
}

// GetTemperature busca a temperatura de um CEP
func (h *TemperatureHandler) GetTemperature(c *gin.Context) {
	cep := c.Param("cep")

	// Validar CEP
	if !h.cepService.ValidateCEP(cep) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "invalid zipcode",
		})
		return
	}

	// Buscar localização do CEP
	location, err := h.cepService.GetLocation(cep)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "can not find zipcode",
		})
		return
	}

	// Buscar temperatura
	temperature, err := h.weatherService.GetTemperature(location.Localidade, location.UF)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "erro ao consultar temperatura",
		})
		return
	}

	// Converter temperaturas
	fahrenheit, kelvin := h.temperatureService.ConvertTemperatures(temperature)

	// Retornar resposta
	response := models.TemperatureResponse{
		TempC: temperature,
		TempF: fahrenheit,
		TempK: kelvin,
	}

	c.JSON(http.StatusOK, response)
}
