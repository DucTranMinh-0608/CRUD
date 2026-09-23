package repositories

import (
	"backend/internal/models"
	"backend/internal/utils"
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

	params := map[string]interface{}{
		"p_id_nguoi_tao": req.IDNguoiTao,
		"p_id_san_pham":  req.IDSanPham,
		"p_nhiem_vu":     req.NhiemVu,
		"p_so_luong":     req.SoLuong,
		"p_ghi_chu":      req.GhiChu,
	}

	data := r.client.Rpc(
		"create_transport_with_stock",
		"exact",
		params,
	)

	if strings.Contains(data, "IDNguoiTao không được null") {
		return nil, utils.ErrEmpty
	}

	if strings.Contains(data, "IDSanPham không được null") {
		return nil, utils.ErrEmpty
	}

	if strings.Contains(data, "SoLuong phải lớn hơn 0") {
		return nil, utils.ErrInvalidQuantity
	}

	if strings.Contains(data, "NhiemVu chỉ được là import hoặc export") {
		return nil, utils.ErrInvalidTask
	}

	if strings.Contains(data, "Không tìm thấy hoặc tài khoản người tạo đã bị vô hiệu hóa") {
		return nil, utils.ErrAccountDisabled
	}

	if strings.Contains(data, "Không tìm thấy sản phẩm ID") {
		return nil, utils.ErrProductDisabled
	}

	if strings.Contains(data, "Không đủ tồn kho") {
		return nil, utils.ErrNegativeStock
	}

	var transports []models.Transport

	if err := json.Unmarshal([]byte(data), &transports); err != nil {
		return nil, fmt.Errorf(
			"failed to decode created transport: %w",
			err,
		)
	}

	if len(transports) == 0 {
		return nil, errors.New(
			"no transport record returned",
		)
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
