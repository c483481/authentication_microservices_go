package main

import (
	"fmt"
	"log"
	"mailer-services/mail"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/goccy/go-json"
)

const (
	appName = "Mailer Service"
	appVersion = "0.0.0"
	appPort = 80
	rpcPort = 5000
	gRpcPort = 50000
)

var mailer Mail


func main() {
	log.Println("Loading .env file...")
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, attempting to use environment variables")
	}

	log.Println("Creating mailer...")
	mailer = createMail()
	log.Println("Mailer created")

	log.Println("Setting up validator...")
	SetupValidate()
	log.Println("Validator setuped")

	err = rpc.Register(&RPCServer{})

	if err != nil {
		log.Fatalf("Error registering RPC server: %v", err)
	}

	go startRPCServer()

	go startGRPCServer()
	
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

	log.Println("Setting routes...")
	app.Post("/mail/send", HandleSendEmail)
	
	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusNotFound).JSON(map[string]any{
			"message": "Not Found",
			"data":    fmt.Sprintf("Not Found path %s", ctx.Path()),
		})
	})

	log.Println("Server started")
	log.Fatal(app.Listen(fmt.Sprintf(":%d", appPort)))
}

func startRPCServer() {
	log.Println("Starting RPC server on port", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", rpcPort))
	if err != nil {
		log.Fatalf("Error starting RPC server: %v", err)
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("Error accepting RPC connection: ", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func startGRPCServer() {
	log.Println("Starting gRPC server on port", gRpcPort)
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", gRpcPort))
	if err != nil {
		log.Fatalf("Error starting gRPC server: %v", err)
	}

	s := grpc.NewServer()
	
	mail.RegisterMailServiceServer(s, &GRPCServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error starting gRPC server: %v", err)
	}
}

