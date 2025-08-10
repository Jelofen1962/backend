// backend/internal/handler/catalog_handler.go
package handler

import (
	"backend/internal/service"
	"backend/pkg/jsonutil"
	"net/http"
	"strconv"
)

type CatalogHandler struct {
	catalogService *service.CatalogService
}

func NewCatalogHandler(s *service.CatalogService) *CatalogHandler {
	return &CatalogHandler{catalogService: s}
}

func (h *CatalogHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters from query string
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20 // Default and max limit
	}

	products, err := h.catalogService.ListProducts(r.Context(), page, limit)
	if err != nil {
		jsonutil.RespondWithError(w, http.StatusInternalServerError, "Could not retrieve products")
		return
	}

	jsonutil.RespondWithJSON(w, http.StatusOK, products)
}

func (h *CatalogHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.catalogService.ListCategories(r.Context())
	if err != nil {
		jsonutil.RespondWithError(w, http.StatusInternalServerError, "Could not retrieve categories")
		return
	}

	jsonutil.RespondWithJSON(w, http.StatusOK, categories)
}
