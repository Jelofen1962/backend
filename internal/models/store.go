// backend/internal/models/store.go
package models

import "time"

// CartItem corresponds to the "cart_items" table.
// This is the raw data structure.
type CartItem struct {
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CartItemDetail is a DTO (Data Transfer Object) used for API responses.
// It enriches the CartItem with details from the products table.
type CartItemDetail struct {
	ProductID     string    `json:"product_id"`
	Quantity      int       `json:"quantity"`
	ProductName   string    `json:"product_name"`   // From products table
	PricePerUnit  float64   `json:"price_per_unit"` // From products table
	LineItemTotal float64   `json:"line_item_total"`
	AddedAt       time.Time `json:"added_at"`
}
