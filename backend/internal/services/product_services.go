package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/utils"
	"context"
	"strings"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *models.CreateProduct) (*models.Product, error)
	GetAllProducts(ctx context.Context, id int64) ([]models.Product, int64, error)
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
		return nil, utils.ErrEmpty
	}
	req.TenMay = strings.TrimSpace(req.TenMay)
	req.Hang = strings.TrimSpace(req.Hang)
	req.TrangThai = strings.TrimSpace(req.TrangThai)

	if req.TenMay == "" {
		return nil, utils.ErrEmpty
	}
	if req.Hang == "" {
		return nil, utils.ErrEmpty
	}
	if req.SoLuong < 0 {
		return nil, utils.ErrNegative
	}
	if req.TrangThai == "" {
		return nil, utils.ErrEmpty
	}
	if req.TrangThai != "enabled" && req.TrangThai != "disabled" {
		return nil, utils.ErrInvalidStatus
	}

	product, err := s.repo.CreateProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) GetAllProducts(ctx context.Context, id int64) ([]models.Product, int64, error) {
	products, total, err := s.repo.GetAllProducts(ctx, id)
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (s *productService) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	if id <= 0 {
		return nil, utils.ErrProductNotFound
	}
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id int64, req *models.UpdateProduct) (*models.Product, error) {
	if id <= 0 {
		return nil, utils.ErrProductNotFound
	}
	if req == nil {
		return nil, utils.ErrEmpty
	}

	if req.TenMay != nil {
		*req.TenMay = strings.TrimSpace(*req.TenMay)
		if *req.TenMay == "" {
			return nil, utils.ErrEmpty
		}
	}

	if req.Hang != nil {
		*req.Hang = strings.TrimSpace(*req.Hang)
		if *req.Hang == "" {
			return nil, utils.ErrEmpty
		}
	}

	if req.SoLuong != nil && *req.SoLuong < 0 {
		return nil, utils.ErrNegative
	}

	if req.TrangThai != nil {
		*req.TrangThai = strings.TrimSpace(*req.TrangThai)
		if *req.TrangThai == "" {
			return nil, utils.ErrEmpty
		}
		if *req.TrangThai != "enabled" && *req.TrangThai != "disabled" {
			return nil, utils.ErrInvalidStatus
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
// 		return utils.ErrProductNotFound
// 	}
// 	return s.repo.DeleteProduct(ctx, id)
// }
