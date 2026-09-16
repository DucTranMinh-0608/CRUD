package controllers

import (
	"errors"
	//"net/http"
	"strconv"

	"go.uber.org/zap"

	"backend/internal/dto"
	"backend/internal/logger"
	"backend/internal/mappers"
	"backend/internal/repositories"
	"backend/internal/services"
	"backend/internal/utils"

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
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		logger.Log.Warn(
			"Dữ liệu gửi lên không hợp lệ",
			zap.Error(err),
		)

		utils.InvalidDataResponse(c, err)
		return
	}

	productModel := mappers.ToProduct(&req)

	product, err := ctrl.service.CreateProduct(c.Request.Context(), productModel)
	if err != nil {
		if errors.Is(err, services.ErrEmptyProductName) ||
			errors.Is(err, services.ErrEmptyBrand) ||
			errors.Is(err, services.ErrNegativeQuantity) {

			logger.Log.Warn(
				"Dữ liệu tạo sản phẩm không hợp lệ",
				zap.Error(err),
			)

			utils.InvalidDataResponse(c, err)
			return
		}

		logger.Log.Error(
			"Lỗi hệ thống khi tạo sản phẩm",
			zap.Error(err),
		)

		utils.NoConnectDataResponse(c, err)
		return
	}

	logger.Log.Info(
		"Tạo sản phẩm thành công",
		zap.Any("product_id", product.ID),
	)

	utils.SuccessResponse(c, "Tạo sản phẩm", product)
}

func (ctrl *ProductController) GetAllProducts(c *gin.Context) {
	products, err := ctrl.service.GetAllProducts(c.Request.Context())
	if err != nil {

		logger.Log.Error(
			"Lỗi hệ thống khi lấy danh sách sản phẩm",
			zap.Error(err),
		)

		utils.NoConnectDataResponse(c, err)
		return
	}

	logger.Log.Info(
		"Lấy danh sách sản phẩm thành công",
		zap.Any("count", len(products)),
	)

	utils.SuccessResponse(c, "Lấy sản phẩm", products)
}

func (ctrl *ProductController) GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {

		logger.Log.Warn(
			"ID sản phẩm không hợp lệ",
			zap.Error(err),
		)

		utils.InvalidDataResponse(c, err)
		return
	}

	product, err := ctrl.service.GetProductByID(c.Request.Context(), id)
	if err != nil {

		if errors.Is(err, repositories.ErrProductNotFound) {

			logger.Log.Warn(
				"Không tìm thấy sản phẩm với ID đã cung cấp",
				zap.Int64("id", id),
				zap.Error(err),
			)

			utils.NotFoundDataResponse(c, err)
			return
		}

		logger.Log.Error(
			"Lỗi hệ thống khi tìm kiếm sản phẩm",
			zap.Error(err),
		)

		utils.NoConnectDataResponse(c, err)
		return
	}

	logger.Log.Info(
		"Tìm thấy sản phẩm với ID đã cung cấp",
		zap.Int64("id", id),
	)

	utils.SuccessResponse(c, "Lấy sản phẩm", product)
}

func (ctrl *ProductController) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {

		logger.Log.Warn(
			"ID sản phẩm không hợp lệ",
			zap.Error(err),
		)

		utils.InvalidDataResponse(c, err)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		logger.Log.Warn(
			"Dữ liệu cập nhật không hợp lệ",
			zap.Error(err),
		)

		utils.InvalidDataResponse(c, err)
		return
	}

	updatedProductModel := mappers.ToUpdateProduct(&req)

	updatedProduct, err := ctrl.service.UpdateProduct(c.Request.Context(), id, updatedProductModel)
	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {

			logger.Log.Warn(
				"Không tìm thấy sản phẩm để cập nhật",
				zap.Int64("id", id),
				zap.Error(err),
			)

			utils.NotFoundDataResponse(c, err)
			return
		}

		if errors.Is(err, services.ErrEmptyProductName) ||
			errors.Is(err, services.ErrEmptyBrand) ||
			errors.Is(err, services.ErrNegativeQuantity) {

			logger.Log.Warn(
				"Dữ liệu cập nhật sản phẩm không hợp lệ",
				zap.Error(err),
			)

			utils.InvalidDataResponse(c, err)
			return
		}

		logger.Log.Error(
			"Lỗi hệ thống khi cập nhật sản phẩm",
			zap.Error(err),
		)

		utils.NoConnectDataResponse(c, err)
		return
	}

	logger.Log.Info(
		"Cập nhật sản phẩm thành công",
		zap.Int64("id", id),
	)

	utils.SuccessResponse(c, "Cập nhật sản phẩm", updatedProduct)
}

func (ctrl *ProductController) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {

		logger.Log.Warn(
			"ID sản phẩm không hợp lệ",
			zap.Error(err),
		)

		utils.InvalidDataResponse(c, err)
		return
	}

	if err := ctrl.service.DeleteProduct(c.Request.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {

			logger.Log.Warn(
				"Không tìm thấy sản phẩm để xóa",
				//zap.Int64("id", id),
				zap.Error(err),
			)

			utils.NotFoundDataResponse(c, err)
			return
		}

		logger.Log.Error(
			"Lỗi hệ thống khi xóa sản phẩm",
			zap.Error(err),
		)

		utils.NoConnectDataResponse(c, err)
		return
	}

	logger.Log.Info(
		"Xóa sản phẩm thành công",
		zap.Int64("id", id),
	)

	utils.SuccessResponse(c, "Xoá sản phẩm", nil)
}
