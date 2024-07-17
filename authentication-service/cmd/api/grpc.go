package main

import (
	"authentication-service/auth"
	"authentication-service/data"
	"context"
	"errors"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCServer struct {
	auth.UnimplementedAuthServiceServer
	Users *data.Users
}

func (s *GRPCServer) Auth(ctx context.Context, req *auth.AuthRequest) (*auth.AuthResponse, error) {
	payload := AuthPayload{
		Email: req.GetEmail(),
		Password: req.GetPassword(),
	}

	if err := ValidateStruct(&payload); err != nil {
		return nil, errors.New(err.Message)
	}

	user, err := s.Users.GetByEmail(payload.Email)

	if err != nil {
		return nil, err
	}

	if ok, err := user.PasswordMatches(payload.Password); err != nil || !ok {
		return nil, errors.New("invalid credentials")
	}

	return &auth.AuthResponse{
		User: &auth.User{
			Id: user.ID,
			Email: user.Email,
			FirstName: user.FirstName,
			LastName: user.LastName,
			Active: user.Active,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}, nil
}
