package repositories

import (
	"backend/internal/models"
	"context"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, req *models.CreateProduct) (*models.Product, error)
	GetAllProducts(ctx context.Context, id int64) ([]models.Product, int64, error)
	GetProductByID(ctx context.Context, id int64) (*models.Product, error)
	UpdateProduct(ctx context.Context, id int64, req *models.UpdateProduct) (*models.Product, error)
	//DeleteProduct(ctx context.Context, id int64) error
}
