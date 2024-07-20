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
	case "log":
		logsPayload := LogsPayload{
			Name: payload.Data["name"].(string),
			Data: payload.Data["data"].(string),
		}
		if valid := ValidateStruct(logsPayload); valid != nil {
			return c.Status(fiber.StatusBadRequest).JSON(valid)
		}

		return insertLogsRPC(c, logsPayload)
	case "mail":
		mailPayload := MailPayload{
			From: payload.Data["from"].(string),
			To: payload.Data["to"].(string),
			Subject: payload.Data["subject"].(string),
			Message: payload.Data["message"].(string),
		}
		if valid := ValidateStruct(mailPayload); valid != nil {
			return c.Status(fiber.StatusBadRequest).JSON(valid)
		}

		return sendMailRPC(c, mailPayload)
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

func insertLogsRPC(c *fiber.Ctx, payload LogsPayload) error {
	client, err := RPCPoolLogs.Get()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
			Success: false,
			Code:    "E_REQ",
			Message: err.Error(),
		})
	}

	var result string

	err = client.Call("RPCServer.LogItems", payload, &result)

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

func sendMailRPC(c *fiber.Ctx, payload MailPayload) error {
	client, err := RPCPoolMail.Get()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponseType{
			Success: false,
			Code:    "E_REQ",
			Message: err.Error(),
		})
	}

	var result string

	err = client.Call("RPCServer.SendMain", payload, &result)

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
