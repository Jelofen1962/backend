// backend/internal/handler/user_handler.go
package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/service"
	"backend/pkg/jsonutil"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

// CreateUser handles the POST /api/v1/users request.
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req service.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Call the business logic layer
	user, err := h.userService.Create(r.Context(), req)
	if err != nil {
		// In a real app, you'd check the error type to return different statuses
		jsonutil.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonutil.RespondWithJSON(w, http.StatusCreated, user)
}

// GetUserByID handles the GET /api/v1/users/{id} request.
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		jsonutil.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	jsonutil.RespondWithJSON(w, http.StatusOK, user)
}
