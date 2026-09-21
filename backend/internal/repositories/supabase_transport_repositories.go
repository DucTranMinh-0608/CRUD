package repositories

import (
	"backend/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type SupabaseTransportRepository struct {
	client *supabase.Client
}

func NewSupabaseTransportRepository(client *supabase.Client) *SupabaseTransportRepository {
	return &SupabaseTransportRepository{
		client: client,
	}
}

func (r *SupabaseTransportRepository) CreateTransport(ctx context.Context, req *models.CreateTransport) (*models.Transport, error) {
	userIdStr := strconv.FormatInt(req.IDNguoiTao, 10)
	productIdStr := strconv.FormatInt(req.IDSanPham, 10)

	if req.TenNguoiTao == "" {
		userData, _, err := r.client.From("users").
			Select("Ten", "", false).
			Eq("ID", userIdStr).
			ExecuteWithContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch user info: %w", err)
		}

		var users []struct {
			Ten string `json:"Ten"`
		}
		if err := json.Unmarshal(userData, &users); err != nil {
			return nil, fmt.Errorf("failed to decode user info: %w", err)
		}
		if len(users) == 0 {
			return nil, ErrUserNotFound
		}
		req.TenNguoiTao = users[0].Ten
	}

	productData, _, err := r.client.From("products").
		Select("TenMay, SoLuong", "", false).
		Eq("ID", productIdStr).
		ExecuteWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current product info: %w", err)
	}

	var products []struct {
		TenMay  string `json:"TenMay"`
		SoLuong int    `json:"SoLuong"`
	}
	if err := json.Unmarshal(productData, &products); err != nil {
		return nil, fmt.Errorf("failed to decode product info: %w", err)
	}
	if len(products) == 0 {
		return nil, ErrProductNotFound
	}

	if req.TenSanPham == "" {
		req.TenSanPham = products[0].TenMay
	}

	var newSoLuong int
	switch strings.ToLower(strings.TrimSpace(req.NhiemVu)) {
	case "import":
		newSoLuong = products[0].SoLuong + int(req.SoLuong)
	case "export":
		newSoLuong = products[0].SoLuong - int(req.SoLuong)
	default:
		return nil, fmt.Errorf("nhiệm vụ không hợp lệ: %s", req.NhiemVu)
	}

	if newSoLuong < 0 {
		return nil, ErrNegativeStock
	}

	data, _, err := r.client.From("transports").
		Insert(req, false, "", "representation", "").
		ExecuteWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transport: %w", err)
	}

	var transports []models.Transport
	if err := json.Unmarshal(data, &transports); err != nil {
		return nil, fmt.Errorf("failed to decode created transport: %w", err)
	}

	if len(transports) == 0 {
		return nil, errors.New("no transport record returned after creation")
	}

	_, _, err = r.client.From("products").
		Update(map[string]interface{}{
			"SoLuong": newSoLuong,
		}, "representation", "").
		Eq("ID", productIdStr).
		ExecuteWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("tạo phiếu vận chuyển thành công nhưng cập nhật tồn kho thất bại: %w", err)
	}

	return &transports[0], nil
}

func (r *SupabaseTransportRepository) GetAllTransports(ctx context.Context) ([]models.Transport, error) {
	data, _, err := r.client.From("transports").
		Select("*", "", false).
		Order("ID", &postgrest.OrderOpts{Ascending: false}).
		ExecuteWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transports: %w", err)
	}

	var transports []models.Transport
	if err := json.Unmarshal(data, &transports); err != nil {
		return nil, fmt.Errorf("failed to decode transports: %w", err)
	}

	return transports, nil
}

func (r *SupabaseTransportRepository) GetTransportsByUserID(ctx context.Context, userID int64) ([]models.Transport, error) {
	userIdStr := strconv.FormatInt(userID, 10)
	data, _, err := r.client.From("transports").
		Select("*", "", false).
		Eq("IDNguoiTao", userIdStr).
		Order("ID", &postgrest.OrderOpts{Ascending: false}).
		ExecuteWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transports by user id: %w", err)
	}

	var transports []models.Transport
	if err := json.Unmarshal(data, &transports); err != nil {
		return nil, fmt.Errorf("failed to decode transports: %w", err)
	}

	return transports, nil
}
