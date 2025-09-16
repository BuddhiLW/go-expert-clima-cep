package handlers

import (
	"net/http"
	"time"

	"cep-temperatura/internal/models"
	"cep-temperatura/internal/services"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	ctx := c.Request.Context()
	tracer := otel.Tracer("service-b")
	
	cep := c.Param("cep")

	// Criar span para validação do CEP
	ctx, span := tracer.Start(ctx, "validate-cep")
	span.SetAttributes(attribute.String("cep", cep))

	// Validar CEP
	if !h.cepService.ValidateCEP(cep) {
		span.SetAttributes(attribute.String("error", "invalid_zipcode"))
		span.End()
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "invalid zipcode",
		})
		return
	}
	span.End()

	// Criar span para busca de localização
	ctx, span = tracer.Start(ctx, "fetch-location")
	span.SetAttributes(attribute.String("cep", cep))

	// Buscar localização do CEP
	location, err := h.cepService.GetLocation(cep)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", "location_not_found"))
		span.End()
		c.JSON(http.StatusNotFound, gin.H{
			"message": "can not find zipcode",
		})
		return
	}

	span.SetAttributes(
		attribute.String("location.city", location.Localidade),
		attribute.String("location.state", location.UF),
	)
	span.End()

	// Criar span para busca de temperatura
	ctx, span = tracer.Start(ctx, "fetch-temperature")
	span.SetAttributes(
		attribute.String("city", location.Localidade),
		attribute.String("state", location.UF),
	)

	start := time.Now()
	// Buscar temperatura
	temperature, err := h.weatherService.GetTemperature(location.Localidade, location.UF)
	duration := time.Since(start)
	
	span.SetAttributes(attribute.Int64("weather.duration_ms", duration.Milliseconds()))

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", "weather_api_error"))
		span.End()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "erro ao consultar temperatura",
		})
		return
	}

	span.SetAttributes(attribute.Float64("temperature.celsius", temperature))
	span.End()

	// Criar span para conversão de temperaturas
	ctx, span = tracer.Start(ctx, "convert-temperatures")
	
	// Converter temperaturas
	fahrenheit, kelvin := h.temperatureService.ConvertTemperatures(temperature)

	span.SetAttributes(
		attribute.Float64("temperature.celsius", temperature),
		attribute.Float64("temperature.fahrenheit", fahrenheit),
		attribute.Float64("temperature.kelvin", kelvin),
	)
	span.End()

	// Retornar resposta
	response := models.TemperatureResponse{
		City:  location.Localidade,
		TempC: temperature,
		TempF: fahrenheit,
		TempK: kelvin,
	}

	c.JSON(http.StatusOK, response)
}
