package main

import (
	"errors"

	"github.com/goccy/go-json"
)

type RPCServer struct{}

func (r *RPCServer) SendMain(payload RequestPayload, resp *string) error {
	if err := ValidateStruct(&payload); err != nil {
		jsonResponse, _ := json.Marshal(&err)
		return errors.New(string(jsonResponse))
	}

	msg := Message {
		From: payload.From,
		To: payload.To,
		Subject: payload.Subject,
		Data: payload.Message,
	}

	err := mailer.SendSMTPMessage(msg)

	if err != nil {
		jsonResponse, _ := json.Marshal(&ErrorResponseType{
			Success: false,
			Code:    "E_ERR",
			Message: err.Error(),
		})
		return errors.New(string(jsonResponse))
	}

	result := &AppResponses{
		Success: true,
		Code: "OK",
		Data: "Sent to "+ payload.To,
	}

	jsonResponse, _ := json.Marshal(result)

	*resp = string(jsonResponse)

	return nil

}
