package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"strings"
)

type TransportService interface {
	CreateTransport(ctx context.Context, req *models.CreateTransport) (*models.Transport, error)
	GetAllTransports(ctx context.Context) ([]models.Transport, error)
	GetTransportsByUserID(ctx context.Context, currentUser *models.User, userID int64) ([]models.Transport, error)
}

type transportService struct {
	repo        repositories.TransportRepository
	userRepo    repositories.UserRepository
	productRepo repositories.ProductRepository
}

func NewTransportService(
	repo repositories.TransportRepository,
	userRepo repositories.UserRepository,
	productRepo repositories.ProductRepository,
) TransportService {
	return &transportService{
		repo:        repo,
		userRepo:    userRepo,
		productRepo: productRepo,
	}
}

func (s *transportService) CreateTransport(ctx context.Context, req *models.CreateTransport) (*models.Transport, error) {
	if req.IDNguoiTao <= 0 {
		return nil, ErrEmpty
	}
	if req.IDSanPham <= 0 {
		return nil, ErrEmpty
	}
	if req.SoLuong <= 0 {
		return nil, ErrInvalidQuantity
	}

	req.NhiemVu = strings.ToLower(strings.TrimSpace(req.NhiemVu))
	if req.NhiemVu == "" {
		return nil, ErrEmpty
	}
	if req.NhiemVu != "import" && req.NhiemVu != "export" {
		return nil, ErrInvalidTask
	}

	user, err := s.userRepo.GetUserByID(ctx, req.IDNguoiTao)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			return nil, repositories.ErrUserNotFound
		}
		return nil, err
	}

	userStatus := strings.ToLower(strings.TrimSpace(user.TrangThai))
	if userStatus == "disabled" {
		return nil, ErrAccountDisabled
	}

	product, err := s.productRepo.GetProductByID(ctx, req.IDSanPham)
	if err != nil {
		if err == repositories.ErrProductNotFound {
			return nil, repositories.ErrProductNotFound
		}
		return nil, err
	}

	productStatus := strings.ToLower(strings.TrimSpace(product.TrangThai))
	if productStatus == "disabled" {
		return nil, ErrProductDisabled
	}

	req.TenNguoiTao = user.Ten
	req.TenSanPham = product.TenMay

	result, err := s.repo.CreateTransport(ctx, req)
	if err != nil {
		return nil, err
	}

	if result.TenNguoiTao == "" {
		result.TenNguoiTao = user.Ten
	}
	if result.TenSanPham == "" {
		result.TenSanPham = product.TenMay
	}

	return result, nil
}

func (s *transportService) GetAllTransports(ctx context.Context) ([]models.Transport, error) {
	transports, err := s.repo.GetAllTransports(ctx)
	if err != nil {
		return nil, err
	}
	return transports, nil
}

func (s *transportService) GetTransportsByUserID(ctx context.Context, currentUser *models.User, userID int64) ([]models.Transport, error) {
	if currentUser == nil {
		return nil, ErrMissingToken
	}

	if userID <= 0 {
		return nil, repositories.ErrUserNotFound
	}

	isAdmin := strings.EqualFold(strings.TrimSpace(currentUser.ChucVu), "admin")
	if !isAdmin && currentUser.ID != userID {
		return nil, ErrForbidden
	}

	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			return nil, repositories.ErrUserNotFound
		}
		return nil, err
	}

	transports, err := s.repo.GetTransportsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return transports, nil
}
