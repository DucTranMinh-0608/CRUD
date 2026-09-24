package repositories

import (
	"backend/internal/models"
	"context"
	"time"
)

type AuthRepository interface {
	Register(ctx context.Context, req *models.Register) (*models.User, error)
	Login(ctx context.Context, req *models.Login) (*models.User, error)
	//GetUserFromToken(ctx context.Context, token string) (*models.User, error)
	CreateJTI(ctx context.Context, jti string, expiresAt time.Time) error
	FindJTI(ctx context.Context, jti string) (bool, error)
	DeleteJTI(ctx context.Context, jti string) error
}
