package main

import (
	"authentication-service/data"

	"github.com/gofiber/fiber/v2"
)

func setAuthRoutes(app fiber.Router, model *data.Users) {
	app.Post("/", func(c *fiber.Ctx) error {
		req := &struct {
			Email    string `json:"email" validate:"required,email,max=255"`
			Password string `json:"password" validate:"required,min=8,max=255"`
		}{}
		
		valid := ValidateBody(c, req)

		if valid != nil {
			return c.Status(valid.Status).JSON(valid)
		}

		user, err := model.GetByEmail(req.Email)

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(&ErrorResponseType{
				Success: false,
				Code:    "E_NOT_FOUND",
				Message: "User not found",
			})
		}

		if ok, err := user.PasswordMatches(req.Password); err != nil || !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(&ErrorResponseType{
				Success: false,
				Code:    "E_UNAUTHORIZED",
				Message: "Invalid credentials",
			})
		}



		return WrapData(c, user)
	})
}