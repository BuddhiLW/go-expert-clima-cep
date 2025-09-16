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

	// Criar handler
	handler := handlers.NewTemperatureHandler(cepService, weatherService, temperatureService)

	// Configurar roteador
	router := gin.Default()
	
	// Middleware de tracing
	router.Use(telemetry.GinMiddleware())
	
	// Rotas
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "service-b"})
	})
	router.GET("/temperature/:cep", handler.GetTemperature)

	// Iniciar servidor
	address := cfg.GetServerAddress()
	log.Printf("Serviço B iniciado em %s", address)
	log.Fatal(router.Run(address))
}
