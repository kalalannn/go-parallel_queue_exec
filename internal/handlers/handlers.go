package handlers

import (
	"go-parallel_queue/internal/messages"
	"log"

	"github.com/gofiber/fiber/v2"
)

func StatusHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"running": []string{"aa", "bb"},
		"queue":   []string{"cc", "dd"},
	})
}

func PlanHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	var data map[string]int
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": messages.InvalidJSON,
		})
	}
	log.Println(data)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": messages.OK,
	})
}
