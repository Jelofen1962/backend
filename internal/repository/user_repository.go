// backend/internal/repository/user_repository.go
package repository

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/models"
)

// UserRepository is an interface that abstracts database operations for users.
// This allows us to easily mock it for testing the service layer.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

// postgresUserRepository is the PostgreSQL implementation of the UserRepository.
type postgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (full_name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, user.FullName, user.Email, user.PasswordHash).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
	)
	return err
}

func (r *postgresUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, full_name, email, password_hash, created_at FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *postgresUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	// Implementation would be very similar to FindByID
	// ...
	return nil, nil // Placeholder
}
