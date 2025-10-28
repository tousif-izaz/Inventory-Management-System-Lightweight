package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"inventory-management/pkg/domain"
	"net/http"
	"os"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var tokenString string

		// Try to get token from Authorization header first
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			// Fallback to cookie
			cookie, err := c.Cookie("token")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "missing or invalid token"})
			}
			tokenString = cookie.Value
		}

		jwtKey := os.Getenv("JWT_KEY")

		claims := &domain.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid token"})
		}

		c.Set("user", claims)
		return next(c)
	}
}
