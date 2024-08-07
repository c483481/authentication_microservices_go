package main

import (
	"broker-service/internal/connection"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/goccy/go-json"
)

const (
	appName = "Broker Service"
	appVersion = "0.0.0"
	appPort = 80
)

var RPCPoolAuth *connection.RPCPool
var RPCPoolLogs *connection.RPCPool
var GRPCPoolAuth *connection.GRPCPool
var GRPCPoolLogs *connection.GRPCPool
var RPCPoolMail *connection.RPCPool
var GRPCPoolMail *connection.GRPCPool

func main() {
	log.Println("Initializing validator...")

	SetupValidate()
	
	log.Println("Successfully initialized validator.")

	RPCPoolAuth = connection.NewRPCPool("authentication-services:5000")
	GRPCPoolAuth = connection.NewGRPCPool("authentication-services:50000")

	RPCPoolLogs = connection.NewRPCPool("logger-services:5000")
	GRPCPoolLogs = connection.NewGRPCPool("logger-services:50000")

	RPCPoolMail = connection.NewRPCPool("mail-services:5000")
	GRPCPoolMail = connection.NewGRPCPool("mail-services:50000")

	log.Println("Starting server...")
	// set up config app
	app := fiber.New(fiber.Config{
		AppName:       appName,
		CaseSensitive: true,             // enable case-sensitive routing
		IdleTimeout:   10 * time.Second, // set idle timeout to 10 second
		ReadTimeout:   5 * time.Second,  // set read timeout to 5 second
		WriteTimeout:  5 * time.Second,  // set write timeout to 5 second
		BodyLimit:     50 * 1024,        // set limit body to 50 KB
		JSONEncoder:   json.Marshal,     // set json encoder to goccy encoder
		JSONDecoder:   json.Unmarshal,   // set json decoder to goccy decoder
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(500).JSON(map[string]string{
				"error": err.Error(),
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, HEAD, PUT, PATCH, POST, DELETE, OPTIONS",
	}))

	// set logger
	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02T15:04:05",
		TimeZone:   "Asia/Jakarta",
	}))

	// set helmet
	app.Use(helmet.New())

	// set recover for keep app alive when there an error from handler
	app.Use(recover.New())

	app.Get("/monitor", monitor.New(monitor.Config{
		Title: appName,
	}))

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(map[string]string{
			"app name": appName,
			"version":  appVersion,
		})
	})

	app.Post("/http", HandleHttpSubmission)
	app.Post("/rpc", HandleRPCSubmission)
	app.Post("/grpc", HandleGRPCSubmission)

	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusNotFound).JSON(map[string]any{
			"message": "Not Found",
			"data":    fmt.Sprintf("Not Found path %s", ctx.Path()),
		})
	})

	log.Println("Server started")
	log.Fatal(app.Listen(fmt.Sprintf(":%d", appPort)))
}