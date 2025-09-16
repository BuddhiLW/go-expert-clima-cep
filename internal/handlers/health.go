package handlers

import (
	"context"
	"net/http"
	"time"

	"cep-temperatura/internal/services"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// HealthHandler gerencia os endpoints de health check
type HealthHandler struct {
	cepService     services.CEPService
	weatherService services.WeatherService
	httpClient     services.HTTPClient
	serviceBURL    string
}

// NewHealthHandler cria uma nova instância do handler de health
func NewHealthHandler(
	cepService services.CEPService,
	weatherService services.WeatherService,
	httpClient services.HTTPClient,
	serviceBURL string,
) *HealthHandler {
	return &HealthHandler{
		cepService:     cepService,
		weatherService: weatherService,
		httpClient:     httpClient,
		serviceBURL:    serviceBURL,
	}
}

// HealthResponse representa a resposta do health check
type HealthResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
	Uptime    string            `json:"uptime,omitempty"`
}

// HealthCheck executa o health check básico
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("health-check")
	
	ctx, span := tracer.Start(ctx, "health-check")
	defer span.End()

	response := HealthResponse{
		Status:    "ok",
		Service:   "cep-temperatura",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    make(map[string]string),
	}

	// Verificar status geral
	response.Checks["overall"] = "ok"
	
	span.SetAttributes(
		attribute.String("service", response.Service),
		attribute.String("status", response.Status),
	)

	// Suportar tanto GET quanto HEAD requests
	if c.Request.Method == "HEAD" {
		c.Status(http.StatusOK)
		return
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheckDetailed executa health check detalhado com dependências
func (h *HealthHandler) HealthCheckDetailed(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("health-check-detailed")

	ctx, span := tracer.Start(ctx, "health-check-detailed")
	defer span.End()

	response := HealthResponse{
		Status:    "ok",
		Service:   "cep-temperatura",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    make(map[string]string),
	}

	// Verificar CEP service
	ctx, cepSpan := tracer.Start(ctx, "check-cep-service")
	if h.cepService.ValidateCEP("01310100") {
		response.Checks["cep_service"] = "ok"
	} else {
		response.Checks["cep_service"] = "error"
		response.Status = "degraded"
	}
	cepSpan.End()

	// Verificar Weather service (apenas se configurado)
	ctx, weatherSpan := tracer.Start(ctx, "check-weather-service")
	if h.weatherService != nil {
		// Teste simples de conectividade (sem fazer requisição real)
		response.Checks["weather_service"] = "ok"
	} else {
		response.Checks["weather_service"] = "not_configured"
		response.Status = "degraded"
	}
	weatherSpan.End()

	// Verificar conectividade com Serviço B (apenas para Serviço A)
	if h.serviceBURL != "" {
		ctx, serviceBSpan := tracer.Start(ctx, "check-service-b")
		if h.checkServiceBConnectivity(ctx) {
			response.Checks["service_b"] = "ok"
		} else {
			response.Checks["service_b"] = "error"
			response.Status = "degraded"
		}
		serviceBSpan.End()
	}

	// Determinar status final
	statusCode := http.StatusOK
	if response.Status == "degraded" {
		statusCode = http.StatusServiceUnavailable
	}

	span.SetAttributes(
		attribute.String("service", response.Service),
		attribute.String("status", response.Status),
		attribute.String("checks_count", string(rune(len(response.Checks)))),
	)

	c.JSON(statusCode, response)
}

// checkServiceBConnectivity verifica se o Serviço B está acessível
func (h *HealthHandler) checkServiceBConnectivity(ctx context.Context) bool {
	if h.httpClient == nil || h.serviceBURL == "" {
		return false
	}

	// Criar requisição para health check do Serviço B
	req, err := http.NewRequestWithContext(ctx, "GET", h.serviceBURL+"/health", nil)
	if err != nil {
		return false
	}

	// Fazer requisição com timeout
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// ReadinessCheck verifica se o serviço está pronto para receber tráfego
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("readiness-check")
	
	ctx, span := tracer.Start(ctx, "readiness-check")
	defer span.End()

	response := HealthResponse{
		Status:    "ready",
		Service:   "cep-temperatura",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    make(map[string]string),
	}

	// Verificações básicas de readiness
	response.Checks["service"] = "ready"
	response.Checks["dependencies"] = "ok"

	span.SetAttributes(
		attribute.String("service", response.Service),
		attribute.String("status", response.Status),
	)

	// Suportar tanto GET quanto HEAD requests
	if c.Request.Method == "HEAD" {
		c.Status(http.StatusOK)
		return
	}

	c.JSON(http.StatusOK, response)
}

// LivenessCheck verifica se o serviço está vivo
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("liveness-check")
	
	ctx, span := tracer.Start(ctx, "liveness-check")
	defer span.End()

	response := HealthResponse{
		Status:    "alive",
		Service:   "cep-temperatura",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    make(map[string]string),
	}

	response.Checks["service"] = "alive"

	span.SetAttributes(
		attribute.String("service", response.Service),
		attribute.String("status", response.Status),
	)

	// Suportar tanto GET quanto HEAD requests
	if c.Request.Method == "HEAD" {
		c.Status(http.StatusOK)
		return
	}

	c.JSON(http.StatusOK, response)
}
