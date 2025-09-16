package main

import (
	"log"
	"os"

	"cep-temperatura/internal/handlers"
	"cep-temperatura/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Configurar Gin para produção
	gin.SetMode(gin.ReleaseMode)

	// Criar instâncias dos serviços
	cepService := services.NewCEPService()
	weatherService := services.NewWeatherService()
	temperatureService := services.NewTemperatureService()

	// Criar handler
	handler := handlers.NewTemperatureHandler(cepService, weatherService, temperatureService)

	// Configurar roteador
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/temperature/:cep", handler.GetTemperature)

	// Obter porta do ambiente ou usar padrão
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado na porta %s", port)
	log.Fatal(router.Run(":" + port))
}
