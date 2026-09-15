package repositories

import (
	"backend/models"
	"context"
	"errors"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error)
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	GetProductByID(ctx context.Context, id int64) (*models.Product, error)
	UpdateProduct(ctx context.Context, id int64, req *models.UpdateProductRequest) (*models.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}
