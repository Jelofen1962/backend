// backend/internal/repository/product_repository.go
package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"errors"
)

// ProductRepository abstracts database operations for products and categories.
type ProductRepository interface {
	FindAll(ctx context.Context, limit, offset int) ([]*models.Product, error)
	FindByID(ctx context.Context, id string) (*models.Product, error)
	// We'll add category methods here as well
	FindAllCategories(ctx context.Context) ([]*models.Category, error)
}

type postgresProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) ProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) FindAll(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	query := `
		SELECT id, category_id, name, description, price, inventory_count, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		p := new(models.Product)
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.InventoryCount, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *postgresProductRepository) FindByID(ctx context.Context, id string) (*models.Product, error) {
	p := new(models.Product)
	query := `
		SELECT id, category_id, name, description, price, inventory_count, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.InventoryCount, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	return p, err
}

func (r *postgresProductRepository) FindAllCategories(ctx context.Context) ([]*models.Category, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM categories ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		c := new(models.Category)
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
