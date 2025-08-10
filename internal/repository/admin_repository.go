// backend/internal/repository/admin_repository.go
package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
)

// AdminRepository abstracts privileged write operations.
type AdminRepository interface {
	CreateProduct(ctx context.Context, product *models.Product) error
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id string) error
	AdjustProductInventory(ctx context.Context, id string, change int) (int, error)
}

type postgresAdminRepository struct {
	db *sql.DB
}

func NewPostgresAdminRepository(db *sql.DB) AdminRepository {
	return &postgresAdminRepository{db: db}
}

func (r *postgresAdminRepository) CreateProduct(ctx context.Context, p *models.Product) error {
	query := `
		INSERT INTO products (category_id, name, description, price, inventory_count)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	// Note: In a real app, you'd insert all the i18n fields too.
	return r.db.QueryRowContext(ctx, query, p.CategoryID, p.Name, p.Description, p.Price, p.InventoryCount).Scan(
		&p.ID, &p.CreatedAt, &p.UpdatedAt,
	)
}

func (r *postgresAdminRepository) UpdateProduct(ctx context.Context, p *models.Product) error {
	query := `
		UPDATE products
		SET category_id = $2, name = $3, description = $4, price = $5, inventory_count = $6, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.CategoryID, p.Name, p.Description, p.Price, p.InventoryCount)
	return err
}

func (r *postgresAdminRepository) DeleteProduct(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresAdminRepository) AdjustProductInventory(ctx context.Context, id string, change int) (int, error) {
	query := `
		UPDATE products
		SET inventory_count = inventory_count + $1, updated_at = NOW()
		WHERE id = $2
		RETURNING inventory_count
	`
	var newInventory int
	err := r.db.QueryRowContext(ctx, query, change, id).Scan(&newInventory)
	return newInventory, err
}
