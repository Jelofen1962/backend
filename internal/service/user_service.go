// backend/internal/service/user_service.go
package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"

	// We'd use a real password hashing library here
	"golang.org/x/crypto/bcrypt"
)

// DTOs (Data Transfer Objects) for requests and responses
type CreateUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	// Business logic: check if user exists
	if _, err := s.repo.FindByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("email already in use")
	}

	// Business logic: hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Prepare the model for the database
	user := &models.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	// Call the database layer
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Map the internal model to a public-facing response (omitting the password hash)
	response := &UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	}
	return response, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err // Let handler decide on 404
	}

	response := &UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	}
	return response, nil
}
