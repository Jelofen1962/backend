// amin-n-co-api/internal/handler/http.go
package handler

import (
	"context"
	"log"
	"net/http"
	"runtime/debug" // Required for printing stack traces

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// =============================================================================
// AUTH MIDDLEWARE (Placeholders)
// In a real application, you would implement the logic in these functions,
// likely calling your pkg/auth/jwt.go functions.
// =============================================================================

// contextKey is a custom type to avoid key collisions in the context.
type contextKey string

const (
	userContextKey = contextKey("userID")
	roleContextKey = contextKey("userRole")
)

// Authenticator is a middleware that verifies a user's JWT token.
// If the token is valid, it extracts the user ID and role and injects them
// into the request's context for later handlers to use.
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// --- Production Logic would go here ---
		// 1. Get token from "Authorization: Bearer <token>" header.
		// 2. Call auth.ValidateToken(tokenString).
		// 3. If valid, extract claims. If not, return 401 Unauthorized.
		// ------------------------------------

		// For now, we inject a MOCK user for testing purposes.
		log.Println("[Auth Middleware] MOCK: Injecting user 'mock-user-id-123' with role 'admin'")
		ctx := context.WithValue(r.Context(), userContextKey, "mock-user-id-123")
		ctx = context.WithValue(ctx, roleContextKey, "admin") // Change to "customer" to test admin routes

		// Call the next handler in the chain with the new, enriched context.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnly is an authorization middleware that checks for the 'admin' role.
// It MUST run *after* the Authenticator has run.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(roleContextKey).(string)

		if !ok || role != "admin" {
			log.Printf("[Auth Middleware] FORBIDDEN: User does not have 'admin' role.")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// User is an admin, proceed to the actual handler.
		next.ServeHTTP(w, r)
	})
}

// =============================================================================
// MAIN ROUTER SETUP
// =============================================================================

// NewRouter creates and configures the main application router.
// It accepts all the handlers for dependency injection.
func NewRouter(
	userHandler *UserHandler,
	catalogHandler *CatalogHandler,
	storeHandler *StoreHandler,
	adminHandler *AdminHandler,
) http.Handler {
	r := chi.NewRouter()

	// --- Standard Middleware ---

	// Recoverer catches panics and prevents the server from crashing.
	// This custom implementation also logs the stack trace for debugging.
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					log.Printf("--- PANIC RECOVERED ---")
					log.Printf("Panic: %+v", rvr)
					// Print the stack trace to the console
					log.Printf("Stack trace:\n%s", debug.Stack())
					// Return a generic 500 error to the client
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	})

	r.Use(middleware.RequestID) // Injects a request ID into the context of each request.
	r.Use(middleware.RealIP)    // Sets a http.Request's RemoteAddr to either X-Real-IP or X-Forwarded-For.
	r.Use(middleware.Logger)    // Logs the start and end of each request with structured data.
	r.Use(middleware.Heartbeat("/ping")) // A health-check endpoint.

	// --- API Route Grouping ---
	r.Route("/api/v1", func(r chi.Router) {

		// == Group 1: Public Routes (No Auth Required) ==
		r.Post("/users/register", userHandler.CreateUser)
		// r.Post("/users/login", userHandler.Login) // Your future login endpoint
		r.Get("/catalog/products", catalogHandler.ListProducts)
		r.Get("/catalog/categories", catalogHandler.ListCategories)
		// r.Get("/catalog/products/{id}", catalogHandler.GetProductByID)

		// == Group 2: Authenticated Routes (User must be logged in) ==
		r.Group(func(r chi.Router) {
			r.Use(Authenticator) // This middleware protects all routes inside this group.

			// User-specific profile routes
			// r.Get("/users/me", userHandler.GetMyProfile)

			// Store routes (cart, orders, wishlist)
			r.Route("/store", func(r chi.Router) {
				r.Get("/cart", storeHandler.GetCart)
				r.Post("/cart/items", storeHandler.AddToCart)
				r.Delete("/cart/items/{productID}", storeHandler.RemoveFromCart)
			})
		})

		// == Group 3: Admin Routes (User must be logged in AND have 'admin' role) ==
		r.Group(func(r chi.Router) {
			r.Use(Authenticator) // First, verify they are a valid user.
			r.Use(AdminOnly)     // Then, check if they have admin permissions.

			r.Route("/admin", func(r chi.Router) {
				r.Post("/products", adminHandler.CreateProduct)
				r.Patch("/products/{id}/inventory", adminHandler.AdjustInventory)
				// Add other admin routes like PUT and DELETE for products here.
			})
		})
	})

	return r
}