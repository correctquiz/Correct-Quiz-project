package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt/v5"
)

var Store *session.Store

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Token"})
		}

		tokenString := ""
		if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		}

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Token Format"})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("my_super_secret_key_12345"), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or Expired Token"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if userIDFloat, ok := claims["user_id"].(float64); ok {
				c.Locals("user_id", uint(userIDFloat))
			}
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Claims"})
		}

		return c.Next()
	}
}
