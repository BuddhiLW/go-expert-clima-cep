package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cep-temperatura/internal/models"
	"cep-temperatura/internal/services"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

// ServiceAHandler gerencia as requisições do Serviço A
type ServiceAHandler struct {
	cepService  services.CEPService
	httpClient  services.HTTPClient
	serviceBURL string
}

// NewServiceAHandler cria uma nova instância do handler do Serviço A
func NewServiceAHandler(
	cepService services.CEPService,
	httpClient services.HTTPClient,
	serviceBURL string,
) *ServiceAHandler {
	return &ServiceAHandler{
		cepService:  cepService,
		httpClient:  httpClient,
		serviceBURL: serviceBURL,
	}
}

// CEPRequest representa a requisição de CEP
type CEPRequest struct {
	CEP string `json:"cep" binding:"required"`
}

// ValidateAndForwardCEP valida o CEP e encaminha para o Serviço B
func (h *ServiceAHandler) ValidateAndForwardCEP(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("service-a")
	
	// Criar span para validação do CEP
	ctx, span := tracer.Start(ctx, "validate-cep")
	defer span.End()

	var req CEPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", "invalid_request_format"))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request format",
		})
		return
	}

	// Validar CEP
	if !h.cepService.ValidateCEP(req.CEP) {
		span.SetAttributes(
			attribute.String("cep", req.CEP),
			attribute.String("error", "invalid_zipcode"),
		)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "invalid zipcode",
		})
		return
	}

	span.SetAttributes(
		attribute.String("cep", req.CEP),
		attribute.Bool("valid", true),
	)

	// Encaminhar para Serviço B
	response, err := h.forwardToServiceB(ctx, req.CEP)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", "service_b_error"))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "erro ao consultar temperatura",
		})
		return
	}

	span.SetAttributes(attribute.String("service_b_response", "success"))
	c.JSON(http.StatusOK, response)
}

// forwardToServiceB encaminha a requisição para o Serviço B
func (h *ServiceAHandler) forwardToServiceB(ctx context.Context, cep string) (*models.TemperatureResponse, error) {
	tracer := otel.Tracer("service-a")
	
	// Criar span para chamada ao Serviço B
	ctx, span := tracer.Start(ctx, "call-service-b")
	defer span.End()

	span.SetAttributes(
		attribute.String("service_b_url", h.serviceBURL),
		attribute.String("cep", cep),
	)

	// Fazer requisição para o Serviço B
	url := fmt.Sprintf("%s/temperature/%s", h.serviceBURL, cep)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Propagar trace context
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	start := time.Now()
	resp, err := h.httpClient.Do(req)
	duration := time.Since(start)
	
	span.SetAttributes(
		attribute.Int64("http.duration_ms", duration.Milliseconds()),
		attribute.Int("http.status_code", resp.StatusCode),
	)

	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("erro na requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		span.SetAttributes(attribute.String("error", "service_b_error"))
		return nil, fmt.Errorf("serviço B retornou status %d", resp.StatusCode)
	}

	var tempResponse models.TemperatureResponse
	if err := json.NewDecoder(resp.Body).Decode(&tempResponse); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	span.SetAttributes(
		attribute.String("response.city", tempResponse.City),
		attribute.Float64("response.temp_c", tempResponse.TempC),
	)

	return &tempResponse, nil
}
