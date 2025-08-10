// backend/internal/service/admin_service.go
package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
)

// DTO for creating a product
type CreateProductRequest struct {
	CategoryID     *string `json:"category_id"`
	Name           string  `json:"name"`
	Description    *string `json:"description"`
	Price          float64 `json:"price"`
	InventoryCount int     `json:"inventory_count"`
}

// DTO for adjusting inventory
type AdjustInventoryRequest struct {
	Change int `json:"change"` // e.g., +10 or -5
}

type AdminService struct {
	adminRepo repository.AdminRepository
	// May need other repos to validate data, e.g., does category_id exist?
	// productRepo repository.ProductRepository
}

func NewAdminService(ar repository.AdminRepository) *AdminService {
	return &AdminService{adminRepo: ar}
}

func (s *AdminService) CreateProduct(ctx context.Context, req CreateProductRequest) (*models.Product, error) {
	// Business validation
	if req.Name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if req.Price < 0 {
		return nil, errors.New("product price cannot be negative")
	}

	product := &models.Product{
		Name:           req.Name,
		Price:          req.Price,
		InventoryCount: req.InventoryCount,
	}
	if req.CategoryID != nil {
		product.CategoryID.String = *req.CategoryID
		product.CategoryID.Valid = true
	}
	if req.Description != nil {
		product.Description.String = *req.Description
		product.Description.Valid = true
	}

	err := s.adminRepo.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *AdminService) AdjustInventory(ctx context.Context, productID string, change int) (int, error) {
	if change == 0 {
		return 0, errors.New("inventory change cannot be zero")
	}
	// Future enhancement: check that inventory doesn't go below zero if change is negative.
	return s.adminRepo.AdjustProductInventory(ctx, productID, change)
}
