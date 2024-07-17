package main

import (
	"broker-service/auth"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func HandleGRPCSubmission(c *fiber.Ctx) error {
	payload := RequestPayload{}

	if valid := ValidateBody(c, &payload); valid != nil {
		return c.Status(fiber.StatusBadRequest).JSON(valid)
	}

	switch payload.Action {
	case "auth":
		authPayload := AuthPayload{
			Email:    payload.Data["email"].(string),
			Password: payload.Data["password"].(string),
		}
		if valid := ValidateStruct(authPayload); valid != nil {
			return c.Status(fiber.StatusBadRequest).JSON(valid)
		}

		return authenticateGRPC(c, &authPayload)
	default:
		return WrapData(c, "Invalid action")
	}
}

func authenticateGRPC(c *fiber.Ctx, payload *AuthPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.NewClient("authentication-services:50000", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
			Success: false,
			Code:    "E_REQ",
			Message: err.Error(),
		})
	}

	defer conn.Close()

	client := auth.NewAuthServiceClient(conn)

	users, err := client.Auth(ctx, &auth.AuthRequest{
		Email: payload.Email,
		Password: payload.Password,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
			Success: false,
			Code:    "E_REQ",
			Message: err.Error(),
		})
	}

	return WrapData(c, users.User)
}