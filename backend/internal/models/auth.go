package models

type Register struct {
	Email       string
	Password    string
	Ten         string
	SoDienThoai string
	ChucVu      string
}

type Login struct {
	Email    string
	Password string
}
