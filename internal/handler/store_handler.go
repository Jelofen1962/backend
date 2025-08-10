// backend/internal/handler/store_handler.go
package handler

import (
	"backend/internal/service"
	"backend/pkg/jsonutil"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type StoreHandler struct {
	storeService *service.StoreService
}

func NewStoreHandler(s *service.StoreService) *StoreHandler {
	return &StoreHandler{storeService: s}
}

// GetCart handles GET /api/v1/store/cart
func (h *StoreHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	// NOTE: We assume a middleware has already authenticated the user
	// and placed their ID in the request context.
	userID := "user-id-from-auth-middleware" // Placeholder

	cart, err := h.storeService.GetCart(r.Context(), userID)
	if err != nil {
		jsonutil.RespondWithError(w, http.StatusInternalServerError, "Could not retrieve cart")
		return
	}
	jsonutil.RespondWithJSON(w, http.StatusOK, cart)
}

// AddToCart handles POST /api/v1/store/cart/items
func (h *StoreHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := "user-id-from-auth-middleware" // Placeholder

	var req service.AddItemToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := h.storeService.AddToCart(r.Context(), userID, req)
	if err != nil {
		jsonutil.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RemoveFromCart handles DELETE /api/v1/store/cart/items/{productID}
func (h *StoreHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID := "user-id-from-auth-middleware" // Placeholder
	productID := chi.URLParam(r, "productID")

	err := h.storeService.RemoveFromCart(r.Context(), userID, productID)
	if err != nil {
		jsonutil.RespondWithError(w, http.StatusInternalServerError, "Could not remove item from cart")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
