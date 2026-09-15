package main

import (
	"net/http"
	"os"

	"go.uber.org/zap"

	"backend/internal/config"
	"backend/internal/logger"
	"backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	logger.Init()
	defer logger.Sync()

	supabaseClient := config.InitSupabaseClient()

	router := gin.Default()
	routes.SetupRoutes(router, supabaseClient)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Log.Info(
		"Server đang lắng nghe kết nối",
		zap.String("port", port),
		zap.String("url", "http://localhost:"+port),
	)

	if err := router.Run(":" + port); err != nil {
		logger.Log.Error("Không thể khởi chạy HTTP server", zap.Error(err))
	}
}
