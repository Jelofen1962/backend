// backend/internal/service/store_service.go
package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
)

type AddItemToCartRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CartResponse struct {
	Items      []*models.CartItemDetail `json:"items"`
	TotalPrice float64                  `json:"total_price"`
	TotalItems int                      `json:"total_items"`
}

type StoreService struct {
	repo repository.StoreRepository
	// We might need other repos in the future, e.g., to check inventory
	// productRepo repository.ProductRepository
}

func NewStoreService(r repository.StoreRepository) *StoreService {
	return &StoreService{repo: r}
}

func (s *StoreService) AddToCart(ctx context.Context, userID string, req AddItemToCartRequest) error {
	if req.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	item := &models.CartItem{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	// Future enhancement: Check if product exists and if there is enough inventory
	// using the productRepo before adding to cart.

	return s.repo.UpsertCartItem(ctx, item)
}

func (s *StoreService) GetCart(ctx context.Context, userID string) (*CartResponse, error) {
	items, err := s.repo.FindCartByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var totalPrice float64
	var totalItems int
	for _, item := range items {
		totalPrice += item.LineItemTotal
		totalItems += item.Quantity
	}

	response := &CartResponse{
		Items:      items,
		TotalPrice: totalPrice,
		TotalItems: totalItems,
	}
	return response, nil
}

func (s *StoreService) RemoveFromCart(ctx context.Context, userID, productID string) error {
	return s.repo.DeleteCartItem(ctx, userID, productID)
}
