package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
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
		fmt.Println(err)
		return BadRequestResponse
	}

	err = validate.Struct(data)

	if err != nil {
		fmt.Println(err)
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

func SendRequest(method, url string, data any) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(req)
}
