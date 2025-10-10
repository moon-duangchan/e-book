package auth

import (
    "errors"
    "strings"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

// RequireAuth returns a Fiber middleware that checks for a valid
// Bearer JWT in the Authorization header. On success, it stores the
// user ID from the `sub` claim in c.Locals("userID").
func RequireAuth() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authz := c.Get("Authorization")
        if authz == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing Authorization header"})
        }

        parts := strings.SplitN(authz, " ", 2)
        if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid Authorization header"})
        }

        secret := getenvDefault("JWT_SECRET", "")
        if secret == "" {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "server misconfigured"})
        }

        tokenStr := strings.TrimSpace(parts[1])

        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, errors.New("unexpected signing method")
            }
            return []byte(secret), nil
        })
        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
        }

        // Extract user ID from `sub` claim if possible
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            if sub, ok := claims["sub"]; ok {
                c.Locals("userID", sub)
            }
        }

        return c.Next()
    }
}

