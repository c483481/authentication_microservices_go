package main

import (
	"logger-services/data"

	"github.com/gofiber/fiber/v2"
)

func setRoutes(app fiber.Router, LogEntry *data.LogEntry) {
	app.Post("/", func(c *fiber.Ctx) error {
		req := &LogRequest{}
		
		if valid := ValidateBody(c, req); valid != nil {
			return c.Status(fiber.StatusBadRequest).JSON(valid)
		}

		entry := &data.LogEntry{
			Name: req.Name,
			Data: req.Data,
		}

		err := LogEntry.Insert(entry)
		
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return WrapData(c, entry)
	})
}
