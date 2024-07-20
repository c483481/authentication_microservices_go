package main

import (
	"context"
	"errors"
	"fmt"
	"mailer-services/mail"
)

type GRPCServer struct {
	mail.UnimplementedMailServiceServer
}

func (s *GRPCServer) SendMail(ctx context.Context, req *mail.MailRequest) (*mail.MailResponse, error) {
	payload := &RequestPayload{
		From: req.GetFrom(),
		To: req.GetTo(),
		Subject: req.GetSubject(),
		Message: req.GetMessage(),
	}

	fmt.Println(payload)

	if err := ValidateStruct(payload); err != nil {
		return nil, errors.New(err.Message)
	}

	msg := Message {
		From: payload.From,
		To: payload.To,
		Subject: payload.Subject,
		Data: payload.Message,
	}

	err := mailer.SendSMTPMessage(msg)
	if err != nil {
		return nil, err
	}

	return &mail.MailResponse{
		Response: "Sent to " + payload.To,
	}, nil

}
