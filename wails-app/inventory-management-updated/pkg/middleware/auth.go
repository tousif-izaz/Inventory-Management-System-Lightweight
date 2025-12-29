package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"ims-intro/pkg/domain"
	"net/http"
	"os"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var tokenString string

		// Try to get token from cookie first
		cookie, err := c.Cookie("token")
		if err == nil && cookie != nil {
			tokenString = cookie.Value
		} else {
			// Fallback to Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "missing or invalid token"})
			}

			// Extract Bearer token
			const bearerPrefix = "Bearer "
			if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid authorization format"})
			}
			tokenString = authHeader[len(bearerPrefix):]
		}

		// Validate JWT token
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

// AdminOnlyMiddleware ensures only admin users can access the route
func AdminOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		if user == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
		}

		claims, ok := user.(*domain.Claims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid user claims"})
		}

		if claims.Role != "admin" {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "admin access required"})
		}

		return next(c)
	}
}
