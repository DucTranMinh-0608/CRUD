package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/utils"
	"context"
	"net/mail"
	"strings"
)

type AuthService interface {
	Register(ctx context.Context, req *models.Register) (*models.User, error)
	Login(ctx context.Context, req *models.Login) (*models.User, string, error)
	GetMe(ctx context.Context, token string) (*models.User, error)
	Logout(ctx context.Context, token string) error
}

type authService struct {
	repo     repositories.AuthRepository
	userrepo repositories.UserRepository
}

func NewAuthService(repo repositories.AuthRepository, userrepo repositories.UserRepository) AuthService {
	return &authService{
		repo:     repo,
		userrepo: userrepo,
	}
}

func (s *authService) Register(ctx context.Context, req *models.Register) (*models.User, error) {
	if req == nil {
		return nil, utils.ErrEmpty
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	req.Ten = strings.TrimSpace(req.Ten)
	req.SoDienThoai = strings.TrimSpace(req.SoDienThoai)
	req.ChucVu = strings.ToLower(strings.TrimSpace(req.ChucVu))

	if req.Email == "" || req.Password == "" || req.Ten == "" || req.SoDienThoai == "" || req.ChucVu == "" {
		return nil, utils.ErrEmpty
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, utils.ErrInvalidEmail
	}

	if len(req.Password) < 6 {
		return nil, utils.ErrPasswordTooShort
	}

	if req.ChucVu != "admin" && req.ChucVu != "staff" {
		return nil, utils.ErrInvalidRole
	}

	user, err := s.repo.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, req *models.Login) (*models.User, string, error) {
	if req == nil {
		return nil, "", utils.ErrEmpty
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		return nil, "", utils.ErrEmpty
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, "", utils.ErrInvalidEmail
	}

	user, err := s.repo.Login(ctx, req)
	if err != nil {
		return nil, "", err
	}

	status := strings.ToLower(strings.TrimSpace(user.TrangThai))
	if status == "disabled" {
		return nil, "", utils.ErrAccountDisabled
	}

	token, jti, expiresAt, err := GenerateToken(user.ID)

	if err != nil {
		return nil, "", utils.ErrGenerateToken
	}

	if err := s.repo.CreateJTI(ctx, jti, expiresAt); err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) GetMe(ctx context.Context, token string) (*models.User, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, utils.ErrMissingToken
	}

	userID, jti, err := VerifyToken(token)

	if err != nil {
		return nil, err
	}

	isValid, err := s.repo.FindJTI(ctx, jti)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, utils.ErrInvalidToken
	}

	user, err := s.userrepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	status := strings.ToLower(strings.TrimSpace(user.TrangThai))
	if status == "disabled" {
		return nil, utils.ErrAccountDisabled
	}

	return user, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return utils.ErrMissingToken
	}

	_, jti, err := VerifyToken(token)
	if err != nil {
		return err
	}

	return s.repo.DeleteJTI(ctx, jti)
}

