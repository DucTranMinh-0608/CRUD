package dto

type RegisterRequest struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"mat_khau" binding:"required"`
	Ten         string `json:"ten" binding:"required"`
	SoDienThoai string `json:"so_dien_thoai" binding:"required"`
	ChucVu      string `json:"chuc_vu" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"mat_khau" binding:"required"`
}
