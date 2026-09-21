package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"strings"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, req *models.UpdateUser) (*models.User, error)
	//DeleteUser(ctx context.Context, id int64) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	if id <= 0 {
		return nil, repositories.ErrUserNotFound
	}
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int64, req *models.UpdateUser) (*models.User, error) {
	if id <= 0 {
		return nil, repositories.ErrUserNotFound
	}
	if req == nil {
		return nil, ErrEmpty
	}

	if req.Ten != nil {
		*req.Ten = strings.TrimSpace(*req.Ten)
		if *req.Ten == "" {
			return nil, ErrEmpty
		}
	}

	if req.SoDienThoai != nil {
		*req.SoDienThoai = strings.TrimSpace(*req.SoDienThoai)
		if *req.SoDienThoai == "" {
			return nil, ErrEmpty
		}
	}

	if req.ChucVu != nil {
		*req.ChucVu = strings.ToLower(strings.TrimSpace(*req.ChucVu))
		if *req.ChucVu != "admin" && *req.ChucVu != "staff" {
			return nil, ErrInvalidRole
		}
	}

	if req.TrangThai != nil && (*req.TrangThai != "enabled" && *req.TrangThai != "disabled") {
		return nil, ErrInvalidStatus
	}

	user, err := s.repo.UpdateUser(ctx, id, req)
	if err != nil {
		return nil, err
	}
	return user, nil
}
