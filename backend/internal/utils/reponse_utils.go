package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InvalidDataResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"Message": "Dữ liệu không hợp lệ",
		"Error":   err.Error(),
	})
}

func NoConnectDataResponse(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"Message": "Không ghi dữ liệu được vào database",
		"Error":   err.Error(),
	})
}

func NotFoundDataResponse(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, gin.H{
		"Message": "Không tìm thấy dữ liệu",
		"Error":   err.Error(),
	})
}

func SuccessResponse(c *gin.Context, s string, products interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"Message": s + " thành công",
		"Data":    products,
	})
}
