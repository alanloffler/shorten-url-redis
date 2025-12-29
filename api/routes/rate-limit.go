package routes

import (
	"strconv"

	"github.com/alanloffler/shorten-url-fiber-redis/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func GetRateLimit(c *fiber.Ctx) error {
	r2 := database.CreateClient(1)
	defer r2.Close()

	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ip":                 c.IP(),
			"rate_limit":         "No limit set",
			"requests_remaining": "Unlimited",
			"reset_in_seconds":   0,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to retrieve rate limit",
		})
	}

	valInt, _ := strconv.Atoi(val)
	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ip":                 c.IP(),
		"requests_remaining": valInt,
		"reset_in_seconds":   int(ttl.Seconds()),
		"reset_in_minutes":   int(ttl.Minutes()),
	})
}

func ClearRateLimit(c *fiber.Ctx) error {
	r2 := database.CreateClient(1)
	defer r2.Close()

	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "No rate limit to clear",
			"ip":      c.IP(),
		})
	}

	err = r2.Del(database.Ctx, c.IP()).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to clear rate limit",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":        "Rate limit cleared successfully",
		"ip":             c.IP(),
		"previous_limit": val,
	})
}
