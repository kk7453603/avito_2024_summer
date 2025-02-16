// Package middlewares provides functionality for handling JWT-based authentication in HTTP requests.
// It includes a Middlewares struct that uses a tokenManager to parse and validate JWT tokens.
package middlewares

import "github.com/golang-jwt/jwt/v5"

// authHeader is the key used to extract the JWT token from the HTTP request header.
const authHeader = "Authorization"

// tokenManager defines the interface for parsing JWT claims.
type tokenManager interface {
	ParseClaims(string) (*jwt.MapClaims, error)
}

// Middlewares provides middleware functionality for handling JWT-based authentication.
type Middlewares struct {
	tknMng tokenManager
}

// NewMiddlewares creates a new instance of Middlewares with the provided tokenManager.
func NewMiddlewares(tokenManager tokenManager) *Middlewares {
	return &Middlewares{tknMng: tokenManager}
}
