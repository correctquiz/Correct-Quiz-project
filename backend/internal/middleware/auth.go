package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store *session.Store

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		var tokenString string

		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			tokenString = c.Cookies("quiz_session")
		}

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		sess, err := Store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "ไม่สามารถดึงข้อมูล session",
			})
		}

		userIDRaw := sess.Get("user_id")
		if userIDRaw == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "ไม่ได้รับอนุญาต: ยังไม่ได้ล็อกอิน",
			})
		}

		var userID uint
		var ok bool

		switch id := userIDRaw.(type) {
		case uint:
			userID = id
			ok = true
		case float64:
			if id > 0 {
				userID = uint(id)
				ok = true
			}
		case int:
			if id > 0 {
				userID = uint(id)
				ok = true
			}
		}

		if !ok || userID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fmt.Sprintf("ไม่ได้รับอนุญาต: User ID เสียหาย (Type: %T)", userIDRaw),
			})
		}
		c.Locals("user_id", userID)

		return c.Next()
	}
}
