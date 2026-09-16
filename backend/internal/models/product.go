package models

type Product struct {
	ID         int64
	TenMay     string
	Hang       string
	MoTa       string
	SoLuong    int
	TrangThai  string
	Created_At string
	Updated_At string
}

type CreateProduct struct {
	TenMay    string
	Hang      string
	SoLuong   int
	MoTa      string
	TrangThai string
}

type UpdateProduct struct {
	TenMay    *string
	Hang      *string
	MoTa      *string
	SoLuong   *int
	TrangThai *string
}
