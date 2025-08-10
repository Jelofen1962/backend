// backend/internal/repository/store_repository.go
package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
)

// StoreRepository abstracts DB operations for cart, orders, etc.
type StoreRepository interface {
	// Cart methods
	UpsertCartItem(ctx context.Context, item *models.CartItem) error
	FindCartByUser(ctx context.Context, userID string) ([]*models.CartItemDetail, error)
	DeleteCartItem(ctx context.Context, userID, productID string) error
	ClearCart(ctx context.Context, userID string) error
}

type postgresStoreRepository struct {
	db *sql.DB
}

func NewPostgresStoreRepository(db *sql.DB) StoreRepository {
	return &postgresStoreRepository{db: db}
}

func (r *postgresStoreRepository) UpsertCartItem(ctx context.Context, item *models.CartItem) error {
	query := `
		INSERT INTO cart_items (user_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, product_id) DO UPDATE
		SET quantity = cart_items.quantity + EXCLUDED.quantity, updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, item.UserID, item.ProductID, item.Quantity)
	return err
}

func (r *postgresStoreRepository) FindCartByUser(ctx context.Context, userID string) ([]*models.CartItemDetail, error) {
	query := `
		SELECT
			ci.product_id,
			ci.quantity,
			p.name,
			p.price,
			ci.created_at
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = $1
		ORDER BY ci.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.CartItemDetail
	for rows.Next() {
		item := new(models.CartItemDetail)
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.ProductName, &item.PricePerUnit, &item.AddedAt); err != nil {
			return nil, err
		}
		item.LineItemTotal = item.PricePerUnit * float64(item.Quantity)
		items = append(items, item)
	}

	return items, nil
}

func (r *postgresStoreRepository) DeleteCartItem(ctx context.Context, userID, productID string) error {
	query := `DELETE FROM cart_items WHERE user_id = $1 AND product_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, productID)
	return err
}

func (r *postgresStoreRepository) ClearCart(ctx context.Context, userID string) error {
	query := `DELETE FROM cart_items WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
