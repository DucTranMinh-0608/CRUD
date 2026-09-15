package repositories

import (
	"backend/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type SupabaseProductRepository struct {
	client *supabase.Client
}

func NewSupabaseProductRepository(client *supabase.Client) *SupabaseProductRepository {
	return &SupabaseProductRepository{
		client: client,
	}
}

func (r *SupabaseProductRepository) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	data, _, err := r.client.From("products").
		Insert(req, false, "", "representation", "").
		ExecuteWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	var products []models.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, fmt.Errorf("failed to decode created product: %w", err)
	}

	if len(products) == 0 {
		return nil, errors.New("no product record returned after creation")
	}

	return &products[0], nil
}

func (r *SupabaseProductRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	data, _, err := r.client.From("products").
		Select("*", "", false).
		Order("id", &postgrest.OrderOpts{Ascending: false}).
		ExecuteWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	var products []models.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, nil
}

func (r *SupabaseProductRepository) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	data, _, err := r.client.From("products").
		Select("*", "", false).
		Eq("id", strconv.FormatInt(id, 10)).
		ExecuteWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product by id: %w", err)
	}

	var products []models.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, fmt.Errorf("failed to decode product: %w", err)
	}

	if len(products) == 0 {
		return nil, ErrProductNotFound
	}

	return &products[0], nil
}

func (r *SupabaseProductRepository) UpdateProduct(ctx context.Context, id int64, req *models.UpdateProductRequest) (*models.Product, error) {
	updates := map[string]interface{}{}

	if req.TenMay != nil {
		updates["ten_may"] = *req.TenMay
	}
	if req.Hang != nil {
		updates["hang"] = *req.Hang
	}
	if req.SoLuong != nil {
		updates["so_luong"] = *req.SoLuong
	}
	if req.MoTa != nil {
		updates["mo_ta"] = *req.MoTa
	}
	if req.TrangThai != nil {
		updates["trang_thai"] = *req.TrangThai
	}

	data, _, err := r.client.From("products").
		Update(updates, "representation", "").
		Eq("id", strconv.FormatInt(id, 10)).
		ExecuteWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	var products []models.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, fmt.Errorf("failed to decode updated product: %w", err)
	}

	if len(products) == 0 {
		return nil, ErrProductNotFound
	}

	return &products[0], nil
}

func (r *SupabaseProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	data, _, err := r.client.From("products").
		Delete("representation", "").
		Eq("id", strconv.FormatInt(id, 10)).
		ExecuteWithContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	var products []models.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return fmt.Errorf("failed to decode delete result: %w", err)
	}

	if len(products) == 0 {
		return ErrProductNotFound
	}

	return nil
}
