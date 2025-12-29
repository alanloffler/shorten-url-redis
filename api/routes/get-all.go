package routes

import (
	"github.com/alanloffler/shorten-url-fiber-redis/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type urlEntry struct {
	Short string `json:"short"`
	URL   string `json:"url"`
}

func GetAll(c *fiber.Ctx) error {
	r := database.CreateClient(0)
	defer r.Close()

	keys, err := r.Keys(database.Ctx, "*").Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to retrieve URLs",
		})
	}

	if len(keys) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"count": 0,
			"urls":  []urlEntry{},
		})
	}

	urls := make([]urlEntry, 0, len(keys))
	for _, key := range keys {
		val, err := r.Get(database.Ctx, key).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			continue
		}

		urls = append(urls, urlEntry{
			Short: key,
			URL:   val,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count": len(urls),
		"urls":  urls,
	})
}
