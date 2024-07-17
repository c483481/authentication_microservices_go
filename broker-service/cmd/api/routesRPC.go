package main

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func HandleRPCSubmission(c *fiber.Ctx) error {
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

		return authenticateRPC(c, authPayload)
	default:
		return WrapData(c, "Invalid action")
	}
}

func authenticateRPC(c *fiber.Ctx, payload AuthPayload) error {
	client, err := RPCPoolAuth.Get()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
			Success: false,
			Code:    "E_REQ",
			Message: err.Error(),
		})
	}
	
	var result string

	err = client.Call("RPCServer.Auth", payload, &result)

	if err != nil {
		var errorResponse ErrorResponseType
    	if err := json.Unmarshal([]byte(err.Error()), &errorResponse); err == nil {
			return c.Status(errorResponse.Status).JSON(&errorResponse)
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
				Success: false,
				Code:    "E_REQ",
				Message: err.Error(),
			})
		}
	}
	var response AppResponses
	
	if err := json.Unmarshal([]byte(result), &response); err  == nil {
		return c.Status(fiber.StatusOK).JSON(&response)
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
			Success: false,
			Code:    "E_REQ",
			Message: err.Error(),
		})
	}
}

