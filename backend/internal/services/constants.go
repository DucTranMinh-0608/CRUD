package services

import "errors"

var (
	ErrEmptyProductName = errors.New("tên máy không được để trống")
	ErrEmptyBrand       = errors.New("hãng không được để trống")
	ErrNegativeQuantity = errors.New("số lượng không được nhỏ hơn 0")
	ErrEmptyStatus      = errors.New("trạng thái không được để trống")
	ErrInvalidStatus    = errors.New("trạng thái không hợp lệ")
)
