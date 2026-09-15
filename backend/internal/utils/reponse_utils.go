package utils

import (
	//"net/http"

	"backend/internal/models"

	"github.com/gin-gonic/gin"
)

func InvalidDataResponse(c *gin.Context, status int, err error) {
	c.JSON(status, models.APIResponse{
		Success: false,
		Message: "Dữ liệu không hợp lệ",
		Error:   err.Error(),
	})
}

func NoConnectDataResponse(c *gin.Context, status int, err error) {
	c.JSON(status, models.APIResponse{
		Success: false,
		Message: "Không ghi dữ liệu được vào database",
		Error:   err.Error(),
	})
}

func SuccessResponse(c *gin.Context, status int, s string, products interface{}) {
	c.JSON(status, models.APIResponse{
		Success: true,
		Message: s + " thành công",
		Data:    products,
	})
}
