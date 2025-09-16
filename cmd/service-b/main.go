package main

import (
	"log"

	"cep-temperatura/internal/config"
	"cep-temperatura/internal/handlers"
	"cep-temperatura/internal/services"
	"cep-temperatura/internal/telemetry"

	"github.com/gin-gonic/gin"
)

func main() {
	// Carregar configuração
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	// Validar configuração
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuração inválida: %v", err)
	}

	// Inicializar telemetria
	shutdown, err := telemetry.InitTracing("service-b", cfg.Telemetry)
	if err != nil {
		log.Fatalf("Erro ao inicializar telemetria: %v", err)
	}
	defer shutdown()

	// Configurar Gin para produção
	gin.SetMode(gin.ReleaseMode)

	// Criar instâncias dos serviços
	cepService := services.NewCEPService()
	weatherService := services.NewWeatherService(cfg)
	temperatureService := services.NewTemperatureService()
	httpClient := services.NewHTTPClient()

	// Criar handlers
	serviceHandler := handlers.NewTemperatureHandler(cepService, weatherService, temperatureService)
	healthHandler := handlers.NewHealthHandler(cepService, weatherService, httpClient, "")

	// Configurar roteador
	router := gin.Default()

	// Middleware de tracing
	router.Use(telemetry.GinMiddleware())

	// Rotas de health
	router.GET("/health", healthHandler.HealthCheck)
	router.HEAD("/health", healthHandler.HealthCheck)
	router.GET("/health/detailed", healthHandler.HealthCheckDetailed)
	router.GET("/ready", healthHandler.ReadinessCheck)
	router.HEAD("/ready", healthHandler.ReadinessCheck)
	router.GET("/live", healthHandler.LivenessCheck)
	router.HEAD("/live", healthHandler.LivenessCheck)

	// Rotas de serviço
	router.GET("/temperature/:cep", serviceHandler.GetTemperature)

	// Iniciar servidor
	address := cfg.GetServerAddress()
	log.Printf("Serviço B iniciado em %s", address)
	log.Fatal(router.Run(address))
}
