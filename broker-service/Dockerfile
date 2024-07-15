FROM golang:1.21-alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLE=0 go build -o api ./cmd/api

RUN chmod +x /app/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api /app/

CMD [ "/app/api" ]
