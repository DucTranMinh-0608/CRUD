package dto

type CreateTransportRequest struct {
	IDSanPham int64  `json:"id_san_pham" binding:"required"`
	NhiemVu   string `json:"nhiem_vu" binding:"required"`
	SoLuong   int64  `json:"so_luong" binding:"required"`
	GhiChu    string `json:"ghi_chu"`
}

type TransportResponse struct {
	ID          int64  `json:"id"`
	IDNguoiTao  int64  `json:"id_nguoi_tao"`
	TenNguoiTao string `json:"ten_nguoi_tao"`
	IDSanPham   int64  `json:"id_san_pham"`
	TenSanPham  string `json:"ten_san_pham"`
	NhiemVu     string `json:"nhiem_vu"`
	SoLuong     int64  `json:"so_luong"`
	GhiChu      string `json:"ghi_chu"`
}
