package main

type RequestPayload struct {
	Action string         `json:"action" validate:"required"`
	Data   map[string]any `json:"data"`
}

type AuthPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LogsPayload struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
	Data string `json:"data" validate:"required,min=3,max=255"`
}

type MailPayload struct {
	From    string `json:"from" validate:"required,email"`
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required,min=3,max=255"`
	Message string `json:"message" validate:"required,min=2,max=65535"`
}
