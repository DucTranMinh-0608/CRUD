package models

type User struct {
	ID          int64
	SupabaseID  string
	Email       string
	Ten         string
	SoDienThoai string
	ChucVu      string
	TrangThai   string
	Created_At  string
	Updated_At  string
}

type UpdateUser struct {
	Ten         *string
	SoDienThoai *string
	ChucVu      *string
	TrangThai   *string
}
