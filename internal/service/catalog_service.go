// backend/internal/service/catalog_service.go
package service

import (
	"backend/internal/repository"
	"context"
)

// ProductResponse is the DTO for a single product sent to the client.
// We can shape this differently from the database model if needed.
type ProductResponse struct {
	ID             string  `json:"id"`
	CategoryID     *string `json:"category_id,omitempty"`
	Name           string  `json:"name"`
	Description    *string `json:"description,omitempty"`
	Price          float64 `json:"price"`
	InventoryCount int     `json:"inventory_count"`
}

// CategoryResponse is the DTO for a category.
type CategoryResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type CatalogService struct {
	repo repository.ProductRepository
}

func NewCatalogService(r repository.ProductRepository) *CatalogService {
	return &CatalogService{repo: r}
}

func (s *CatalogService) ListProducts(ctx context.Context, page, limit int) ([]*ProductResponse, error) {
	offset := (page - 1) * limit
	products, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Map DB models to response DTOs
	var response []*ProductResponse
	for _, p := range products {
		res := &ProductResponse{
			ID:             p.ID,
			Name:           p.Name,
			Price:          p.Price,
			InventoryCount: p.InventoryCount,
		}
		if p.CategoryID.Valid {
			res.CategoryID = &p.CategoryID.String
		}
		if p.Description.Valid {
			res.Description = &p.Description.String
		}
		response = append(response, res)
	}

	return response, nil
}

func (s *CatalogService) ListCategories(ctx context.Context) ([]*CategoryResponse, error) {
	categories, err := s.repo.FindAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	var response []*CategoryResponse
	for _, c := range categories {
		res := &CategoryResponse{ID: c.ID, Name: c.Name}
		if c.Description.Valid {
			res.Description = &c.Description.String
		}
		response = append(response, res)
	}

	return response, nil
}
