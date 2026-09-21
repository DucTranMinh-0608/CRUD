package mappers

import (
	"backend/internal/dto"
	"backend/internal/models"
)

func ToCreateTransport(req *dto.CreateTransportRequest) *models.CreateTransport {
	return &models.CreateTransport{
		IDSanPham: req.IDSanPham,
		NhiemVu:   req.NhiemVu,
		SoLuong:   req.SoLuong,
		GhiChu:    req.GhiChu,
	}
}

func ToTransportResponse(req *models.Transport) *dto.TransportResponse {
	return &dto.TransportResponse{
		ID:          req.ID,
		IDNguoiTao:  req.IDNguoiTao,
		TenNguoiTao: req.TenNguoiTao,
		IDSanPham:   req.IDSanPham,
		TenSanPham:  req.TenSanPham,
		NhiemVu:     req.NhiemVu,
		SoLuong:     req.SoLuong,
		GhiChu:      req.GhiChu,
	}
}
