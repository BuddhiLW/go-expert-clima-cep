package main

import (
	"log"

	"cep-temperatura/internal/config"
	"cep-temperatura/internal/handlers"
	"cep-temperatura/internal/services"

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
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/temperature/:cep", handler.GetTemperature)

	// Iniciar servidor
	address := cfg.GetServerAddress()
	log.Printf("Servidor iniciado em %s", address)
	log.Fatal(router.Run(address))
}
