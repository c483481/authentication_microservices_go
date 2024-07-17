package main

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
)



func HandleHttpSubmission(c *fiber.Ctx) error {
	payload := RequestPayload{}
	
	if valid := ValidateBody(c, &payload); valid != nil {
		return c.Status(fiber.StatusBadRequest).JSON(valid)
	}

	switch payload.Action {
	case "auth":
		authPayload := AuthPayload{
			Email: payload.Data["email"].(string),
			Password: payload.Data["password"].(string),
		}
		if valid := ValidateStruct(authPayload); valid != nil {
			return c.Status(fiber.StatusBadRequest).JSON(valid)
		}

		return authenticate(c, &authPayload)
	case "log":
		logsPayload := &LogsPayload{
			Name: payload.Data["name"].(string),
			Data: payload.Data["data"].(string),
		}
		if valid := ValidateStruct(logsPayload); valid != nil {
			return c.Status(fiber.StatusBadRequest).JSON(valid)
		}

		return insertLogs(c, logsPayload)

	default:
		return WrapData(c, "Invalid action")
	}
}

func authenticate(c *fiber.Ctx, payload *AuthPayload) error {
	response, err := SendRequest("POST", "http://authentication-services/auth", payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
			Success: false,
			Code:    "E_REQ",
			Message: "Failed to create request",
		})
	}
	
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		var jsonResponse AppResponses
		err := json.NewDecoder(response.Body).Decode(&jsonResponse)


		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
				Success: false,
				Code:    "E_REQ",
				Message: "Failed to create request",
			})
		}

		return WrapData(c, jsonResponse.Data)
	case http.StatusUnauthorized:
		return c.Status(fiber.StatusUnauthorized).JSON(&ErrorResponseType{
			Success: false,
			Code:    "E_AUTH",
			Message: "Invalid credentials",
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
			Success: false,
			Code:    "E_REQ",
			Message: "Failed to create request",
		})
	}
}

func insertLogs(c *fiber.Ctx, payload *LogsPayload) error {
	response, err := SendRequest("POST", "http://logger-services/logs", payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
			Success: false,
			Code:    "E_REQ",
			Message: err.Error(),
		})
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
			Success: false,
			Code:    "E_REQ",
			Message: "Failed to create request",
		})
	}

	var jsonResponse AppResponses
	err = json.NewDecoder(response.Body).Decode(&jsonResponse)


	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
			Success: false,
			Code:    "E_REQ",
			Message: "Failed to create request",
		})
	}

	return WrapData(c, jsonResponse.Data)
}
