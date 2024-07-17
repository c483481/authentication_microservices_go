package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func SetupValidate() {
	validate = validator.New()
}

type AppResponses struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Data    any    `json:"data,omitempty"`
}
type ErrorResponseType struct {
	Success bool 	`json:"success"`
	Status  int 	`json:"-"`
	Code    string 	`json:"code"`
	Message string 	`json:"message"`
	Data    any    	`json:"data,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

var (
	BadRequestResponse = &ErrorResponseType{
		Success: false,
		Status:  fiber.StatusBadRequest,
		Code:    "E_BAD_REQUEST",
		Message: "Bad Request",
	}
)


func WrapData(ctx *fiber.Ctx, data interface{}) error {
	return ctx.JSON(&AppResponses{
		Success: true,
		Code:    "OK",
		Data:    data,
	})
}

func ValidateBody(ctx *fiber.Ctx, data any) *ErrorResponseType {
	err := ctx.BodyParser(data)

	if err != nil {
		return BadRequestResponse
	}

	err = validate.Struct(data)

	if err != nil {
		return BadRequestResponse
	}

	return nil
}

func ValidateStruct(data any) *ErrorResponseType {
	err := validate.Struct(data)

	if err != nil {
		return BadRequestResponse
	}

	return nil
}
