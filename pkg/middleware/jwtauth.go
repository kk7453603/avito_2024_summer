package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type JWTMiddlewareConfig struct {
	SecretKey []byte
}

// Я мог бы добавить кастомных или стандартных claims в этот JWT токен, но не вижу в этом смысла из-за отсутствия ролей у пользователей в этом проекте,
func JWTMiddleware(config JWTMiddlewareConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Отсутствует заголовок Authorization",
				})
			}

			//Bearer <token>
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Неверный формат заголовка Authorization",
				})
			}
			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Неподходящий метод подписи")
				}
				return config.SecretKey, nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Неверный или истёкший JWT",
				})
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				c.Set("userID", claims["sub"])
			}

			return next(c)
		}
	}
}

func GenerateToken(config JWTMiddlewareConfig, userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    "avito_api",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(config.SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
