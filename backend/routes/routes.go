package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

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
