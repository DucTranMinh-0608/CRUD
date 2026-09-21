package repositories

import (
	"backend/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

type SupabaseAuthRepository struct {
	client *supabase.Client
}

func NewSupabaseAuthRepository(client *supabase.Client) *SupabaseAuthRepository {
	return &SupabaseAuthRepository{
		client: client,
	}
}

func (r *SupabaseAuthRepository) Register(ctx context.Context, req *models.Register) (*models.User, error) {

	existingData, _, err := r.client.From("users").
		Select("ID", "", false).
		Eq("Email", req.Email).
		ExecuteWithContext(ctx)
	if err == nil && len(existingData) > 0 {
		var existingUsers []struct {
			ID int64 `json:"ID"`
		}
		if json.Unmarshal(existingData, &existingUsers) == nil && len(existingUsers) > 0 {
			return nil, ErrEmailAlreadyExists
		}
	}

	signupReq := types.SignupRequest{
		Email:    req.Email,
		Password: req.Password,
		Data: map[string]interface{}{
			"Ten":         req.Ten,
			"ChucVu":      req.ChucVu,
			"SoDienThoai": req.SoDienThoai,
		},
	}

	resp, err := r.client.Auth.Signup(signupReq)
	if err != nil {
		errLower := strings.ToLower(err.Error())
		if strings.Contains(errLower, "already") || strings.Contains(errLower, "exists") {
			return nil, ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("lỗi khi đăng ký auth supabase: %w", err)
	}

	userId := resp.User.ID
	if userId == uuid.Nil {
		userId = resp.Session.User.ID
	}
	if userId == uuid.Nil {
		return nil, errors.New("không nhận được Supabase ID từ phản hồi auth")
	}

	supabaseID := userId.String()

	data, _, err := r.client.From("users").
		Select("*", "", false).
		Eq("SupabaseID", supabaseID).
		ExecuteWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi lấy thông tin người dùng từ database: %w", err)
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("lỗi khi giải mã thông tin người dùng: %w", err)
	}

	if len(users) == 0 {
		return nil, ErrUserNotFound
	}

	return &users[0], nil
}

func (r *SupabaseAuthRepository) Login(ctx context.Context, req *models.Login) (*models.User, string, error) {
	resp, err := r.client.Auth.SignInWithEmailPassword(req.Email, req.Password)
	if err != nil {
		errLower := strings.ToLower(err.Error())
		if strings.Contains(errLower, "invalid") || strings.Contains(errLower, "credentials") || strings.Contains(errLower, "grant") {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", fmt.Errorf("lỗi khi đăng nhập supabase auth: %w", err)
	}

	userId := resp.User.ID
	if userId == uuid.Nil {
		userId = resp.Session.User.ID
	}
	if userId == uuid.Nil {
		return nil, "", errors.New("không nhận được Supabase ID từ phản hồi auth")
	}

	accessToken := resp.AccessToken
	if accessToken == "" {
		accessToken = resp.Session.AccessToken
	}

	data, _, err := r.client.From("users").
		Select("*", "", false).
		Eq("SupabaseID", userId.String()).
		ExecuteWithContext(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("lỗi khi lấy thông tin người dùng từ database: %w", err)
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, "", fmt.Errorf("lỗi khi giải mã thông tin người dùng: %w", err)
	}

	if len(users) == 0 {
		return nil, "", ErrUserNotFound
	}

	return &users[0], accessToken, nil
}

func (r *SupabaseAuthRepository) GetUserFromToken(ctx context.Context, token string) (*models.User, error) {
	userResp, err := r.client.Auth.WithToken(token).GetUser()
	if err != nil {
		return nil, ErrInvalidToken
	}

	if userResp == nil || userResp.ID == uuid.Nil {
		return nil, ErrInvalidToken
	}

	data, _, err := r.client.From("users").
		Select("*", "", false).
		Eq("SupabaseID", userResp.ID.String()).
		ExecuteWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi lấy thông tin người dùng từ database: %w", err)
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("lỗi khi giải mã thông tin người dùng: %w", err)
	}

	if len(users) == 0 {
		return nil, ErrUserNotFound
	}

	return &users[0], nil
}

func (r *SupabaseAuthRepository) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return r.client.Auth.WithToken(token).Logout()
}
