package repositories

import (
	"backend/internal/models"
	"context"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, req *models.UpdateUser) (*models.User, error)
	//DeleteUser(ctx context.Context, id int64) error
}
