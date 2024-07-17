package main

import (
	"authentication-service/data"
	"errors"

	"github.com/goccy/go-json"
)

type RPCServer struct {
	Users *data.Users
}

func (r *RPCServer) Auth(payload AuthPayload, resp *string) error {
	if err := ValidateStruct(&payload); err != nil {
		jsonResponse, _ := json.Marshal(&err)
		return errors.New(string(jsonResponse))
	}

	user, err := r.Users.GetByEmail(payload.Email)

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
