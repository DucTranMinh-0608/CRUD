package repositories

import (
	"backend/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type SupabaseUserRepository struct {
	client *supabase.Client
}

func NewSupabaseUserRepository(client *supabase.Client) *SupabaseUserRepository {
	return &SupabaseUserRepository{
		client: client,
	}
}

func (r *SupabaseUserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	data, _, err := r.client.From("users").
		Select("*", "", false).
		Order("ID", &postgrest.OrderOpts{Ascending: false}).
		ExecuteWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, nil
}

func (r *SupabaseUserRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	data, _, err := r.client.From("users").
		Select("*", "", false).
		Eq("ID", strconv.FormatInt(id, 10)).
		ExecuteWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user by id: %w", err)
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	if len(users) == 0 {
		return nil, ErrUserNotFound
	}

	return &users[0], nil
}

func (r *SupabaseUserRepository) UpdateUser(ctx context.Context, id int64, req *models.UpdateUser) (*models.User, error) {
	updates := map[string]interface{}{}

	if req.Ten != nil {
		updates["Ten"] = *req.Ten
	}
	if req.SoDienThoai != nil {
		updates["SoDienThoai"] = *req.SoDienThoai
	}
	if req.ChucVu != nil {
		updates["ChucVu"] = *req.ChucVu
	}
	if req.TrangThai != nil {
		updates["TrangThai"] = *req.TrangThai
	}

	if len(updates) == 0 {
		return r.GetUserByID(ctx, id)
	}

	data, _, err := r.client.From("users").
		Update(updates, "representation", "").
		Eq("ID", strconv.FormatInt(id, 10)).
		ExecuteWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to decode updated user: %w", err)
	}

	if len(users) == 0 {
		return nil, ErrUserNotFound
	}

	return &users[0], nil
}
