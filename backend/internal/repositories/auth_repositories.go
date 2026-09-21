package repositories

import (
	"backend/internal/models"
	"context"
	"errors"
)

var (
	ErrEmailAlreadyExists = errors.New("email đã tồn tại")
	ErrInvalidCredentials = errors.New("email hoặc mật khẩu không chính xác")
	ErrInvalidToken       = errors.New("token không hợp lệ hoặc đã hết hạn")
)

type AuthRepository interface {
	Register(ctx context.Context, req *models.Register) (*models.User, error)
	Login(ctx context.Context, req *models.Login) (*models.User, string, error)
	GetUserFromToken(ctx context.Context, token string) (*models.User, error)
	Logout(ctx context.Context, token string) error
}
