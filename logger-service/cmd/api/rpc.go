package main

import (
	"errors"
	"logger-services/data"

	"github.com/goccy/go-json"
)

type RPCServer struct {
	LogEntry *data.LogEntry
}

func (r *RPCServer) LogItems(req *LogRequest, resp *string) error {
	if err := ValidateStruct(&req); err != nil {
		jsonResponse, _ := json.Marshal(&err)
		return errors.New(string(jsonResponse))
	}

	entry := &data.LogEntry{
		Name: req.Name,
		Data: req.Data,
	}

	err := r.LogEntry.Insert(entry)

	if err != nil {
		jsonResponse, _ := json.Marshal(&ErrorResponseType{
			Success: false,
			Code:    "E_CONN",
			Message: err.Error(),
		})
		return errors.New(string(jsonResponse))
	}

	result := &AppResponses{
		Success: true,
		Code: "OK",
		Data: entry,
	}

	jsonResponse, _ := json.Marshal(result)

	*resp = string(jsonResponse)

	return nil
}
