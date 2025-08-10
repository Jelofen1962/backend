// backend/pkg/auth/jwt.go
package auth

import (
	"backend/internal/models" // We might need the user model for role info
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// In a real app, this should be loaded securely from config and be much longer!
var jwtSecretKey = []byte("a-very-secret-and-long-key-for-amin-n-co")

// Claims defines the data stored inside the JWT.
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new signed JWT for a given user.
func GenerateToken(user *models.User) (string, error) {
	// For now, we'll hardcode the role. In a real system, you'd look this up.
	userRole := "customer" // Default role
	// if user.IsAdmin { userRole = "admin" }

	// Set the expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour) // Token is valid for 24 hours

	claims := &Claims{
		UserID: user.ID,
		Role:   userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "backend",
		},
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken parses and validates a JWT string.
// It returns the claims if the token is valid.
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Make sure the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
