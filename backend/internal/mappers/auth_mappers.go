package mappers

import (
	"backend/internal/dto"
	"backend/internal/models"
)

func ToRegister(req *dto.RegisterRequest) *models.Register {
	return &models.Register{
		Email:       req.Email,
		Password:    req.Password,
		Ten:         req.Ten,
		SoDienThoai: req.SoDienThoai,
		ChucVu:      req.ChucVu,
	}
}

func ToLogin(req *dto.LoginRequest) *models.Login {
	return &models.Login{
		Email:    req.Email,
		Password: req.Password,
	}
}
