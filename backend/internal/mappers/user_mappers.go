package mappers

import (
	"backend/internal/dto"
	"backend/internal/models"
)

func ToUpdateUser(req *dto.UpdateUserRequest) *models.UpdateUser {
	return &models.UpdateUser{
		Ten:         req.Ten,
		SoDienThoai: req.SoDienThoai,
		ChucVu:      req.ChucVu,
		TrangThai:   req.TrangThai,
	}
}

func ToUserResponse(req *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:          req.ID,
		Email:       req.Email,
		Ten:         req.Ten,
		SoDienThoai: req.SoDienThoai,
		ChucVu:      req.ChucVu,
		TrangThai:   req.TrangThai,
	}
}
