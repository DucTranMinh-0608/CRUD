package repositories

import (
	"backend/internal/models"
	"context"
	"errors"
)

var (
	ErrTransportNotFound = errors.New("transport not found")
	ErrNegativeStock     = errors.New("số lượng sản phẩm sau khi thực hiện không được âm")
)

type TransportRepository interface {
	CreateTransport(ctx context.Context, req *models.CreateTransport) (*models.Transport, error)
	GetAllTransports(ctx context.Context) ([]models.Transport, error)
	GetTransportsByUserID(ctx context.Context, userID int64) ([]models.Transport, error)
}
