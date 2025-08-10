// backend/internal/handler/admin_handler.go
package handler

import (
	"backend/internal/service"
	"backend/pkg/jsonutil"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(s *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: s}
}

// CreateProduct handles POST /api/v1/admin/products
func (h *AdminHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req service.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	product, err := h.adminService.CreateProduct(r.Context(), req)
	if err != nil {
		jsonutil.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonutil.RespondWithJSON(w, http.StatusCreated, product)
}

// AdjustInventory handles PATCH /api/v1/admin/products/{id}/inventory
func (h *AdminHandler) AdjustInventory(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	var req service.AdjustInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	newStock, err := h.adminService.AdjustInventory(r.Context(), productID, req.Change)
	if err != nil {
		jsonutil.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonutil.RespondWithJSON(w, http.StatusOK, map[string]int{"new_inventory_count": newStock})
}
