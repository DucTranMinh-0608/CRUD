package routes

import (
	"backend/controllers"
	"backend/repositories"
	"backend/services"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

func SetupRoutes(r *gin.Engine, supabaseClient *supabase.Client) {

	api := r.Group("/api")

	productRepo := repositories.NewSupabaseProductRepository(supabaseClient)
	productService := services.NewProductService(productRepo)
	productController := controllers.NewProductController(productService)

	ProductRoutes(api, productController)
}

func ProductRoutes(router *gin.RouterGroup, controller *controllers.ProductController) {
	products := router.Group("/products")
	{
		products.POST("", controller.CreateProduct)
		products.GET("", controller.GetAllProducts)
		products.GET("/:id", controller.GetProductByID)
		products.PUT("/:id", controller.UpdateProduct)
		products.DELETE("/:id", controller.DeleteProduct)
	}
}
