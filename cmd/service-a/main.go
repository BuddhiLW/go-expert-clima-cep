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
	shutdown, err := telemetry.InitTracing("service-a", cfg.Telemetry)
	if err != nil {
		log.Fatalf("Erro ao inicializar telemetria: %v", err)
	}
	defer shutdown()

	// Configurar Gin para produção
	gin.SetMode(gin.ReleaseMode)

	// Criar instâncias dos serviços
	cepService := services.NewCEPService()
	httpClient := services.NewHTTPClient()

	// Criar handlers
	serviceHandler := handlers.NewServiceAHandler(cepService, httpClient, cfg.ServiceB.URL)
	healthHandler := handlers.NewHealthHandler(cepService, nil, httpClient, cfg.ServiceB.URL)

	// Configurar roteador
	router := gin.Default()
	
	// Middleware de tracing
	router.Use(telemetry.GinMiddleware())
	
	// Rotas
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "service-a"})
	})
	router.POST("/cep", handler.ValidateAndForwardCEP)

	// Iniciar servidor
	address := cfg.GetServerAddress()
	log.Printf("Serviço A iniciado em %s", address)
	log.Fatal(router.Run(address))
}
