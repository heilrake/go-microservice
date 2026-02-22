package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const defaultJWTSecret = "ride-sharing-auth-secret-change-in-production"

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"` // "user" | "driver"
	jwt.RegisteredClaims
}

// ParseToken parses and validates a JWT token string
func ParseToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = defaultJWTSecret
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

// ClaimsFromRequest extracts JWT claims from Authorization: Bearer <token> header
func ClaimsFromRequest(r *http.Request) (*Claims, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil, fmt.Errorf("missing Authorization header")
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid Authorization header format")
	}
	return ParseToken(parts[1])
}

// UserIDFromRequest extracts user ID from Authorization: Bearer <token> header
func UserIDFromRequest(r *http.Request) (string, error) {
	claims, err := ClaimsFromRequest(r)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// RoleFromRequest extracts role from JWT in request
func RoleFromRequest(r *http.Request) (string, error) {
	claims, err := ClaimsFromRequest(r)
	if err != nil {
		return "", err
	}
	return claims.Role, nil
}
