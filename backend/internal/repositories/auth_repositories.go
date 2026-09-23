package repositories

import (
	"backend/internal/models"
	"context"
)

type AuthRepository interface {
	Register(ctx context.Context, req *models.Register) (*models.User, error)
	Login(ctx context.Context, req *models.Login) (*models.User, string, error)
	GetUserFromToken(ctx context.Context, token string) (*models.User, error)
	Logout(ctx context.Context, token string) error
}
