package repositories

import (
	"backend/internal/models"
	"context"
)

type TransportRepository interface {
	CreateTransport(ctx context.Context, req *models.CreateTransport) (*models.Transport, error)
	GetAllTransports(ctx context.Context) ([]models.Transport, error)
	GetTransportsByUserID(ctx context.Context, userID int64) ([]models.Transport, error)
}
