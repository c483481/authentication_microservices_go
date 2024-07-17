package main

import (
	"context"
	"errors"
	"logger-services/data"
	"logger-services/logs"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCServer struct {
	logs.UnimplementedLogServiceServer
	Logs *data.LogEntry
}

func (s *GRPCServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	payload := &LogRequest{
		Name: req.GetName(),
		Data: req.GetData(),
	}

	if err := ValidateStruct(payload); err != nil {
		return nil, errors.New(err.Message)
	}

	entry := &data.LogEntry{
		Name: req.Name,
		Data: req.Data,
	}

	err := s.Logs.Insert(entry)

	if err != nil {
		return nil, err
	}

	return &logs.LogResponse{
		Name: entry.Name,
		Data: entry.Data,
		CreatedAt: timestamppb.New(entry.CreatedAt),
		UpdatedAt: timestamppb.New(entry.UpdatedAt),
	}, nil
}
