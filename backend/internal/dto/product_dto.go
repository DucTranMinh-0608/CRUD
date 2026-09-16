package dto

type CreateProductRequest struct {
	TenMay    string `json:"ten_may" binding:"required"`
	Hang      string `json:"hang" binding:"required"`
	SoLuong   int    `json:"so_luong" binding:"required"`
	MoTa      string `json:"mo_ta"`
	TrangThai string `json:"trang_thai" binding:"required"`
}

type UpdateProductRequest struct {
	TenMay    *string `json:"ten_may"`
	Hang      *string `json:"hang"`
	SoLuong   *int    `json:"so_luong"`
	MoTa      *string `json:"mo_ta"`
	TrangThai *string `json:"trang_thai"`
}
