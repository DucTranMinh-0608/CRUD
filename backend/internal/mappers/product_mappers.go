package mappers

import (
	"backend/internal/dto"
	"backend/internal/models"
)

func ToProduct(req *dto.CreateProductRequest) *models.CreateProduct {
	var soLuong int
	if req.SoLuong != nil {
		soLuong = *req.SoLuong
	}
	return &models.CreateProduct{
		TenMay:    req.TenMay,
		Hang:      req.Hang,
		SoLuong:   soLuong,
		MoTa:      req.MoTa,
		TrangThai: req.TrangThai,
	}
}

func ToUpdateProduct(req *dto.UpdateProductRequest) *models.UpdateProduct {
	return &models.UpdateProduct{
		TenMay:    req.TenMay,
		Hang:      req.Hang,
		SoLuong:   req.SoLuong,
		MoTa:      req.MoTa,
		TrangThai: req.TrangThai,
	}
}

func ToProductResponse(req *models.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:        req.ID,
		TenMay:    req.TenMay,
		Hang:      req.Hang,
		MoTa:      req.MoTa,
		SoLuong:   req.SoLuong,
		TrangThai: req.TrangThai,
	}
}

func ToProductListResponse(products []models.Product, total int64, page int64) *dto.ProductListResponse {
	const limit = 10
	responses := make([]dto.ProductResponse, 0, len(products))
	for _, p := range products {
		responses = append(responses, *ToProductResponse(&p))
	}
	return &dto.ProductListResponse{
		Total:    total,
		Page:     page,
		Limit:    limit,
		Products: responses,
	}
}
