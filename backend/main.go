package main

import (
	//"log"
	"net/http"
	"os"

	"go.uber.org/zap"

	"backend/config"
	"backend/controllers"
	"backend/logger"
	"backend/repositories"
	"backend/routes"
	"backend/services"

	"github.com/gin-gonic/gin"
)

func main() {

	logger.Init()
	defer logger.Sync()

	supabaseClient := config.InitSupabaseClient()

	productRepo := repositories.NewSupabaseProductRepository(supabaseClient)
	productService := services.NewProductService(productRepo)
	productController := controllers.NewProductController(productService)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	api := router.Group("/api")
	routes.ProductRoutes(api, productController)

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
