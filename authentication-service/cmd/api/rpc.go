package main

import (
	"authentication-service/data"
	"errors"

	"github.com/goccy/go-json"
)

type RPCServer struct {
	users *data.Users
}

type RPCPayload struct {
	Email string `validate:"required,email,max=255"`
	Password string `validate:"required,min=8,max=255"`
}

func (r *RPCServer) Auth(payload RPCPayload, resp *string) error {
	if err := ValidateStruct(&payload); err != nil {
		jsonResponse, _ := json.Marshal(&payload)
		return errors.New(string(jsonResponse))
	}

	user, err := r.users.GetByEmail(payload.Email)

	if err != nil {
		jsonResponse, _ := json.Marshal(&ErrorResponseType{
			Success: false,
			Code:    "E_NOT_FOUND",
			Message: "User not found",
		})
		return errors.New(string(jsonResponse))
	}

	if ok, err := user.PasswordMatches(payload.Password); err != nil || !ok {
		jsonResponse, _ := json.Marshal(&ErrorResponseType{
			Success: false,
			Code:    "E_UNAUTHORIZED",
			Message: "Invalid credentials",
		})
		return errors.New(string(jsonResponse))
	}

	result := &AppResponses{
		Success: true,
		Code: "OK",
		Data: user,
	}

	jsonResponse, _ := json.Marshal(result)

	*resp = string(jsonResponse)

	return nil
}
