package dto

type UpdateUserRequest struct {
	Ten         *string `json:"ten"`
	SoDienThoai *string `json:"so_dien_thoai"`
	ChucVu      *string `json:"chuc_vu"`
	TrangThai   *string `json:"trang_thai"`
}

type UserResponse struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	Ten         string `json:"ten"`
	SoDienThoai string `json:"so_dien_thoai"`
	ChucVu      string `json:"chuc_vu"`
	TrangThai   string `json:"trang_thai"`
}
