package main

import (
	"broker-service/auth"
	"broker-service/logs"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
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
	case "log":
		logsPayload := &LogsPayload{
			Name: payload.Data["name"].(string),
			Data: payload.Data["data"].(string),
		}
		if valid := ValidateStruct(logsPayload); valid != nil {
			return c.Status(fiber.StatusBadRequest).JSON(valid)
		}

		return insertLogsGRPC(c, logsPayload)
	default:
		return WrapData(c, "Invalid action")
	}
}

func authenticateGRPC(c *fiber.Ctx, payload *AuthPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := GRPCPoolAuth.Get()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
			Success: false,
			Code:    "E_REQ",
			Message: err.Error(),
		})
	}

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

func insertLogsGRPC(c *fiber.Ctx, payload *LogsPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := GRPCPoolLogs.Get()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
			Success: false,
			Code:    "E_REQ",
			Message: err.Error(),
		})
	}

	client := logs.NewLogServiceClient(conn)

	item, err := client.WriteLog(ctx, &logs.LogRequest{
		Name: payload.Name,
		Data: payload.Data,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
			Success: false,
			Code:    "E_REQ",
			Message: err.Error(),
		})
	}

	return WrapData(c, item) 
}
