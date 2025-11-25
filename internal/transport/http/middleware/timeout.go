package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Timeout(duration time.Duration) fiber.Handler {
    return func(c *fiber.Ctx) error {
        ctx, cancel := context.WithTimeout(c.Context(), duration)
        defer cancel()
		
        c.SetUserContext(ctx)

        err := c.Next()

        if ctx.Err() == context.DeadlineExceeded {
            return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
                "error": fiber.Map{
                    "code":    "REQUEST_TIMEOUT",
                    "message": "request executed too long",
                },
            })
        }

        return err
    }
}
