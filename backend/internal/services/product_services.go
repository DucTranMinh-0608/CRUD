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
	//DeleteProduct(ctx context.Context, id int64) error
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
	if req == nil {
		return nil, ErrEmpty
	}
	req.TenMay = strings.TrimSpace(req.TenMay)
	req.Hang = strings.TrimSpace(req.Hang)
	req.TrangThai = strings.TrimSpace(req.TrangThai)

	if req.TenMay == "" {
		return nil, ErrEmpty
	}
	if req.Hang == "" {
		return nil, ErrEmpty
	}
	if req.SoLuong < 0 {
		return nil, ErrNegative
	}
	if req.TrangThai == "" {
		return nil, ErrEmpty
	}
	if req.TrangThai != "enabled" && req.TrangThai != "disabled" {
		return nil, ErrInvalidStatus
	}

	product, err := s.repo.CreateProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	products, err := s.repo.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *productService) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	if id <= 0 {
		return nil, repositories.ErrProductNotFound
	}
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id int64, req *models.UpdateProduct) (*models.Product, error) {
	if id <= 0 {
		return nil, repositories.ErrProductNotFound
	}
	if req == nil {
		return nil, ErrEmpty
	}

	if req.TenMay != nil {
		*req.TenMay = strings.TrimSpace(*req.TenMay)
		if *req.TenMay == "" {
			return nil, ErrEmpty
		}
	}

	if req.Hang != nil {
		*req.Hang = strings.TrimSpace(*req.Hang)
		if *req.Hang == "" {
			return nil, ErrEmpty
		}
	}

	if req.SoLuong != nil && *req.SoLuong < 0 {
		return nil, ErrNegative
	}

	if req.TrangThai != nil {
		*req.TrangThai = strings.TrimSpace(*req.TrangThai)
		if *req.TrangThai == "" {
			return nil, ErrEmpty
		}
		if *req.TrangThai != "enabled" && *req.TrangThai != "disabled" {
			return nil, ErrInvalidStatus
		}
	}

	product, err := s.repo.UpdateProduct(ctx, id, req)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// func (s *productService) DeleteProduct(ctx context.Context, id int64) error {
// 	if id <= 0 {
// 		return repositories.ErrProductNotFound
// 	}
// 	return s.repo.DeleteProduct(ctx, id)
// }
