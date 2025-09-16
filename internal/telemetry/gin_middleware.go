package telemetry

import (
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"github.com/gin-gonic/gin"
)

// GinMiddleware retorna o middleware de tracing para Gin
func GinMiddleware() gin.HandlerFunc {
	return otelgin.Middleware("cep-temperatura")
}
