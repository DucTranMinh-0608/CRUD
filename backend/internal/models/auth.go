package models

import "time"

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

type AccessToken struct {
	JTI       string    `json:"jti"`
	ExpiresAt time.Time `json:"expires_at"`
}
