package utils

import (
	//"net/http"

	"backend/internal/models"

	"github.com/gin-gonic/gin"
)

func InvalidDataResponse(c *gin.Context, err error) {
	c.JSON(400, models.APIResponse{
		Success: false,
		Message: "Dữ liệu không hợp lệ",
		Error:   err.Error(),
	})
}

func NoConnectDataResponse(c *gin.Context, err error) {
	c.JSON(500, models.APIResponse{
		Success: false,
		Message: "Không ghi dữ liệu được vào database",
		Error:   err.Error(),
	})
}

func NotFoundDataResponse(c *gin.Context, err error) {
	c.JSON(404, models.APIResponse{
		Success: false,
		Message: "Không tìm thấy dữ liệu",
		Error:   err.Error(),
	})
}

func SuccessResponse(c *gin.Context, s string, products interface{}) {
	c.JSON(200, models.APIResponse{
		Success: true,
		Message: s + " thành công",
		Data:    products,
	})
}
