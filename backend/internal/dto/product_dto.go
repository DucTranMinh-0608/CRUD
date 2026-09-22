package dto

type CreateProductRequest struct {
	TenMay    string `json:"ten_may" binding:"required"`
	Hang      string `json:"hang" binding:"required"`
	SoLuong   *int   `json:"so_luong" binding:"required,gte=0"`
	MoTa      string `json:"mo_ta"`
	TrangThai string `json:"trang_thai" binding:"required"`
}

type UpdateProductRequest struct {
	TenMay    *string `json:"ten_may"`
	Hang      *string `json:"hang"`
	SoLuong   *int    `json:"so_luong" binding:"omitempty,gte=0"`
	MoTa      *string `json:"mo_ta"`
	TrangThai *string `json:"trang_thai"`
}

type ProductResponse struct {
	ID        int64  `json:"id"`
	TenMay    string `json:"ten_may"`
	Hang      string `json:"hang"`
	MoTa      string `json:"mo_ta"`
	SoLuong   int    `json:"so_luong"`
	TrangThai string `json:"trang_thai"`
}

type ProductListResponse struct {
	Total    int64             `json:"total"`
	Page     int64             `json:"page"`
	Limit    int               `json:"limit"`
	Products []ProductResponse `json:"products"`
}

