package services

import (
	"backend/internal/dto"
	"backend/internal/mappers"
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"strings"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *models.CreateProduct) (*dto.ProductResponse, error)
	GetAllProducts(ctx context.Context) ([]dto.ProductResponse, error)
	GetProductByID(ctx context.Context, id int64) (*dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, id int64, req *models.UpdateProduct) (*dto.ProductResponse, error)
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

func (s *productService) CreateProduct(ctx context.Context, req *models.CreateProduct) (*dto.ProductResponse, error) {
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

	product, err := s.repo.CreateProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	return mappers.ToPorductReponse(product), nil
}

func (s *productService) GetAllProducts(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := s.repo.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ProductResponse, 0, len(products))
	for _, p := range products {
		responses = append(responses, *mappers.ToPorductReponse(&p))
	}
	return responses, nil
}

func (s *productService) GetProductByID(ctx context.Context, id int64) (*dto.ProductResponse, error) {
	if id <= 0 {
		return nil, repositories.ErrProductNotFound
	}
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mappers.ToPorductReponse(product), nil
}

func (s *productService) UpdateProduct(ctx context.Context, id int64, req *models.UpdateProduct) (*dto.ProductResponse, error) {
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

	product, err := s.repo.UpdateProduct(ctx, id, req)
	if err != nil {
		return nil, err
	}
	return mappers.ToPorductReponse(product), nil
}

func (s *productService) DeleteProduct(ctx context.Context, id int64) error {
	if id <= 0 {
		return repositories.ErrProductNotFound
	}
	return s.repo.DeleteProduct(ctx, id)
}
