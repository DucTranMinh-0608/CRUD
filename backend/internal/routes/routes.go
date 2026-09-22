package routes

import (
	"backend/internal/controllers"
	"backend/internal/middlewares"
	"backend/internal/repositories"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

func SetupRoutes(r *gin.Engine, supabaseClient *supabase.Client) {

	api := r.Group("/api")

	authRepo := repositories.NewSupabaseAuthRepository(supabaseClient)
	authService := services.NewAuthService(authRepo)
	authController := controllers.NewAuthController(authService)

	authMiddleware := middlewares.AuthMiddleware(authService)

	AuthRoutes(api, authController, authMiddleware)

	productRepo := repositories.NewSupabaseProductRepository(supabaseClient)
	productService := services.NewProductService(productRepo)
	productController := controllers.NewProductController(productService)
	ProductRoutes(api, productController, authMiddleware)

	userRepo := repositories.NewSupabaseUserRepository(supabaseClient)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)
	UserRoutes(api, userController, authMiddleware)

	transportRepo := repositories.NewSupabaseTransportRepository(supabaseClient)
	transportService := services.NewTransportService(transportRepo, userRepo, productRepo)
	transportController := controllers.NewTransportController(transportService)
	TransportRoutes(api, transportController, authMiddleware)
}

func AuthRoutes(router *gin.RouterGroup, controller *controllers.AuthController, authMiddleware gin.HandlerFunc) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", authMiddleware, middlewares.RequireRoles("admin"), controller.Register)
		auth.POST("/login", controller.Login)
		auth.GET("/me", authMiddleware, controller.GetMe)
		auth.POST("/logout", controller.Logout)
	}
}

func ProductRoutes(router *gin.RouterGroup, controller *controllers.ProductController, authMiddleware gin.HandlerFunc) {
	products := router.Group("/products")
	products.Use(authMiddleware)
	{
		products.GET("/page/:id", controller.GetAllProducts)
		products.GET("/:id", controller.GetProductByID)
		products.POST("", middlewares.RequireRoles("admin"), controller.CreateProduct)
		products.PUT("/:id", middlewares.RequireRoles("admin"), controller.UpdateProduct)
		//products.DELETE("/:id", middlewares.RequireRoles("admin"), controller.DeleteProduct)
	}
}

func UserRoutes(router *gin.RouterGroup, controller *controllers.UserController, authMiddleware gin.HandlerFunc) {
	users := router.Group("/users")
	users.Use(authMiddleware)
	{
		users.GET("", middlewares.RequireRoles("admin"), controller.GetAllUsers)
		users.GET("/:id", middlewares.RequireRoles("admin"), controller.GetUserByID)
		users.PUT("/:id", middlewares.RequireRoles("admin"), controller.UpdateUser)
	}
}

func TransportRoutes(router *gin.RouterGroup, controller *controllers.TransportController, authMiddleware gin.HandlerFunc) {
	transports := router.Group("/transport")
	transports.Use(authMiddleware)
	{
		transports.POST("", middlewares.RequireRoles("staff"), controller.CreateTransport)
		transports.GET("", middlewares.RequireRoles("admin"), controller.GetAllTransports)
		transports.GET("/user/:id", controller.GetTransportsByUserID)
	}
}
