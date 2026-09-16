package mappers

import (
	"backend/internal/dto"
	"backend/internal/models"
)

func ToProduct(req *dto.CreateProductRequest) *models.CreateProduct {
	return &models.CreateProduct{
		TenMay:    req.TenMay,
		Hang:      req.Hang,
		SoLuong:   req.SoLuong,
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
