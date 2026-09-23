package models

type CreateTransport struct {
	IDNguoiTao int64
	IDSanPham  int64
	NhiemVu    string
	SoLuong    int64
	GhiChu     string
}

type Transport struct {
	ID          int64
	IDNguoiTao  int64
	TenNguoiTao string
	IDSanPham   int64
	TenSanPham  string
	NhiemVu     string
	SoLuong     int64
	GhiChu      string
	Created_At  string
	Updated_At  string
}
