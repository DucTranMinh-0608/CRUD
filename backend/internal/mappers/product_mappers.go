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
