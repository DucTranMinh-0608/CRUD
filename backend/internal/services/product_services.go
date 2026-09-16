package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"strings"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *models.CreateProduct) (*models.Product, error)
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	GetProductByID(ctx context.Context, id int64) (*models.Product, error)
	UpdateProduct(ctx context.Context, id int64, req *models.UpdateProduct) (*models.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *models.CreateProduct) (*models.Product, error) {
	req.TenMay = strings.TrimSpace(req.TenMay)
	req.Hang = strings.TrimSpace(req.Hang)
	req.TrangThai = strings.TrimSpace(req.TrangThai)

	if req.TenMay == "" {
		return nil, ErrEmptyProductName
	}
	if req.Hang == "" {
		return nil, ErrEmptyBrand
	}
	if req.SoLuong < 0 {
		return nil, ErrNegativeQuantity
	}
	if req.TrangThai == "" {
		return nil, ErrEmptyStatus
	}
	if req.TrangThai != "enabled" && req.TrangThai != "disabled" {
		return nil, ErrInvalidStatus
	}

	return s.repo.CreateProduct(ctx, req)
}

func (s *productService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetAllProducts(ctx)
}

func (s *productService) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	if id <= 0 {
		return nil, repositories.ErrProductNotFound
	}
	return s.repo.GetProductByID(ctx, id)
}

func (s *productService) UpdateProduct(ctx context.Context, id int64, req *models.UpdateProduct) (*models.Product, error) {
	if id <= 0 {
		return nil, repositories.ErrProductNotFound
	}

	if req.TenMay != nil {
		*req.TenMay = strings.TrimSpace(*req.TenMay)
		if *req.TenMay == "" {
			return nil, ErrEmptyProductName
		}
	}

	if req.Hang != nil {
		*req.Hang = strings.TrimSpace(*req.Hang)
		if *req.Hang == "" {
			return nil, ErrEmptyBrand
		}
	}

	if req.SoLuong != nil && *req.SoLuong < 0 {
		return nil, ErrNegativeQuantity
	}

	if req.TrangThai != nil && (*req.TrangThai != "enabled" && *req.TrangThai != "disabled") {
		return nil, ErrInvalidStatus
	}

	return s.repo.UpdateProduct(ctx, id, req)
}

func (s *productService) DeleteProduct(ctx context.Context, id int64) error {
	if id <= 0 {
		return repositories.ErrProductNotFound
	}
	return s.repo.DeleteProduct(ctx, id)
}
