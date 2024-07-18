package main

import "github.com/gofiber/fiber/v2"

func setUpRoutes(app fiber.Router, mail Mail) {
	app.Post("/send", func(c *fiber.Ctx) error {
		req := &RequestPayload{}

		if valid := ValidateBody(c, req); valid != nil {
			return c.Status(valid.Status).JSON(valid)
		}

		msg := Message {
			From: req.From,
			To: req.To,
			Subject: req.Subject,
			Data: req.Message,
		}

		err := mail.SendSMTPMessage(msg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return WrapData(c, "Sent to "+req.To)
	})
}