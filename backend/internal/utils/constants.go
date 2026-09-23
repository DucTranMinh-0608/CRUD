package utils

import "errors"

var (
	ErrEmpty              = errors.New("không được truyền vào xâu rỗng")
	ErrInvalidRole        = errors.New("chức vụ không hợp lệ, chỉ chấp nhận admin hoặc staff")
	ErrNegative           = errors.New("không được nhỏ hơn 0")
	ErrInvalidQuantity    = errors.New("số phải lớn hơn 0")
	ErrInvalidStatus      = errors.New("trạng thái không hợp lệ, chỉ chấp nhận enabled hoặc disabled")
	ErrInvalidTask        = errors.New("nhiệm vụ không hợp lệ, chỉ chấp nhận import hoặc export")
	ErrNegativeStock      = errors.New("số lượng sản phẩm sau khi thực hiện không được âm")
	ErrUserNotFound       = errors.New("user not found")
	ErrProductNotFound    = errors.New("product not found")
	ErrEmailAlreadyExists = errors.New("email đã tồn tại")
	ErrInvalidCredentials = errors.New("email hoặc mật khẩu không chính xác")
	ErrInvalidToken       = errors.New("token không hợp lệ hoặc đã hết hạn")
	ErrMissingToken       = errors.New("thiếu token xác thực")
	ErrAccountDisabled    = errors.New("tài khoản đã bị vô hiệu hóa")
	ErrProductDisabled    = errors.New("sản phẩm đã bị vô hiệu hóa")
	ErrForbidden          = errors.New("bạn không có quyền thực hiện chức năng này")
	ErrPasswordTooShort   = errors.New("mật khẩu phải có ít nhất 6 ký tự")
	ErrInvalidEmail       = errors.New("định dạng email không hợp lệ")
)
