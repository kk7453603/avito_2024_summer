// Package middlewares provides functionality for handling JWT-based authentication in HTTP requests.
// It includes a Middlewares struct that uses a tokenManager to parse and validate JWT tokens.
package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware is a middleware function that validates JWT tokens in incoming requests.
// It ensures that the request contains a valid "Authorization" header with a Bearer token.
// If the token is valid, it extracts the user ID and username from the token claims and sets them in the context.
func (m *Middlewares) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the "Authorization" header from the request.
		header := c.GetHeader(authHeader)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "the 'Authorization' header is missing"})
			return
		}

		// Split the header into parts to extract the token.
		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "invalid token format"})
			return
		}
		if len(parts[1]) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "token is empty"})
			return
		}

		// Parse and validate the token claims.
		claims, err := m.tknMng.ParseClaims(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
			return
		}

		// Set the user ID and username in the context for use in subsequent handlers.
		c.Set("user_id", (*claims)["sub"])
		c.Set("username", (*claims)["username"])
		c.Next() // Proceed to the next handler.
	}
}
