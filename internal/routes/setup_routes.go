package routes

import (
	"AuthService/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/generate_tokens", handlers.GenerateTokens)

	router.POST("/refresh_tokens", handlers.RefreshTokens)
}
