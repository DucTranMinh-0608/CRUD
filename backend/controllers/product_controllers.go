package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"backend/logger"
	"backend/models"
	"backend/repositories"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service services.ProductService
}

func NewProductController(service services.ProductService) *ProductController {
	return &ProductController{
		service: service,
	}
}

func (ctrl *ProductController) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		logger.Log.Warn(
			"Dữ liệu gửi lên không hợp lệ",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Dữ liệu gửi lên không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	product, err := ctrl.service.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrEmptyProductName) ||
			errors.Is(err, services.ErrEmptyBrand) ||
			errors.Is(err, services.ErrNegativeQuantity) {

			logger.Log.Warn(
				"Dữ liệu tạo sản phẩm không hợp lệ",
				zap.Error(err),
			)

			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "Kiểm tra dữ liệu không đạt yêu cầu",
				Error:   err.Error(),
			})
			return
		}

		logger.Log.Error(
			"Lỗi hệ thống khi tạo sản phẩm",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Lỗi hệ thống khi tạo sản phẩm",
			Error:   err.Error(),
		})
		return
	}

	logger.Log.Info(
		"Tạo sản phẩm thành công",
		zap.Any("product_id", product.ID),
	)

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Tạo sản phẩm thành công",
		Data:    product,
	})
}

func (ctrl *ProductController) GetAllProducts(c *gin.Context) {
	products, err := ctrl.service.GetAllProducts(c.Request.Context())
	if err != nil {

		logger.Log.Error(
			"Lỗi hệ thống khi lấy danh sách sản phẩm",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Lỗi hệ thống khi lấy danh sách sản phẩm",
			Error:   err.Error(),
		})
		return
	}

	logger.Log.Info(
		"Lấy danh sách sản phẩm thành công",
		zap.Any("count", len(products)),
	)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Lấy danh sách sản phẩm thành công",
		Data:    products,
	})
}

func (ctrl *ProductController) GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {

		logger.Log.Warn(
			"ID sản phẩm không hợp lệ",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "ID sản phẩm không hợp lệ",
			Error:   "ID phải là một số nguyên",
		})
		return
	}

	product, err := ctrl.service.GetProductByID(c.Request.Context(), id)
	if err != nil {

		logger.Log.Warn(
			"Không tìm thấy sản phẩm với ID đã cung cấp",
			zap.Int64("id", id),
			zap.Error(err),
		)

		if errors.Is(err, repositories.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Không tìm thấy sản phẩm với ID đã cung cấp",
				Error:   err.Error(),
			})
			return
		}

		logger.Log.Error(
			"Lỗi hệ thống khi tìm kiếm sản phẩm",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Lỗi hệ thống khi tìm kiếm sản phẩm",
			Error:   err.Error(),
		})
		return
	}

	logger.Log.Info(
		"Tìm thấy sản phẩm với ID đã cung cấp",
		zap.Int64("id", id),
	)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Tìm thấy sản phẩm",
		Data:    product,
	})
}

func (ctrl *ProductController) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {

		logger.Log.Warn(
			"ID sản phẩm không hợp lệ",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "ID sản phẩm không hợp lệ",
			Error:   "ID phải là một số nguyên",
		})
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		logger.Log.Warn(
			"Dữ liệu cập nhật không hợp lệ",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Dữ liệu cập nhật không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	updatedProduct, err := ctrl.service.UpdateProduct(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {

			logger.Log.Warn(
				"Không tìm thấy sản phẩm để cập nhật",
				zap.Int64("id", id),
				zap.Error(err),
			)

			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Không tìm thấy sản phẩm để cập nhật",
				Error:   err.Error(),
			})
			return
		}

		if errors.Is(err, services.ErrEmptyProductName) ||
			errors.Is(err, services.ErrEmptyBrand) ||
			errors.Is(err, services.ErrNegativeQuantity) {

			logger.Log.Warn(
				"Dữ liệu cập nhật sản phẩm không hợp lệ",
				zap.Error(err),
			)

			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "Kiểm tra dữ liệu không đạt yêu cầu",
				Error:   err.Error(),
			})
			return
		}

		logger.Log.Error(
			"Lỗi hệ thống khi cập nhật sản phẩm",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Lỗi hệ thống khi cập nhật sản phẩm",
			Error:   err.Error(),
		})
		return
	}

	logger.Log.Info(
		"Cập nhật sản phẩm thành công",
		zap.Int64("id", id),
	)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Cập nhật sản phẩm thành công",
		Data:    updatedProduct,
	})
}

func (ctrl *ProductController) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {

		logger.Log.Warn(
			"ID sản phẩm không hợp lệ",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "ID sản phẩm không hợp lệ",
			Error:   "ID phải là một số nguyên",
		})
		return
	}

	if err := ctrl.service.DeleteProduct(c.Request.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {

			logger.Log.Warn(
				"Không tìm thấy sản phẩm để xóa",
				//zap.Int64("id", id),
				zap.Error(err),
			)

			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Không tìm thấy sản phẩm để xóa",
				Error:   err.Error(),
			})
			return
		}

		logger.Log.Error(
			"Lỗi hệ thống khi xóa sản phẩm",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Lỗi hệ thống khi xóa sản phẩm",
			Error:   err.Error(),
		})
		return
	}

	logger.Log.Info(
		"Xóa sản phẩm thành công",
		zap.Int64("id", id),
	)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Xóa sản phẩm thành công",
	})
}
