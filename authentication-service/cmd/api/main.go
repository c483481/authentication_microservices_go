package main

import (
	"authentication-service/data"
	"authentication-service/migration"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/goccy/go-json"
)

const (
	appName = "Authentication Service"
	appVersion = "0.0.0"
	appPort = 80
)

func main() {
	log.Println("Loading environment variables...")
	err := godotenv.Load()
	
	if err != nil {
		log.Println("Error loading .env file, attempting to use environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	log.Println("Opening database...")

	db, err := openDB(databaseURL)
	
	if err != nil {
		log.Fatal(err)
	}
	
	defer db.Close()

	log.Println("Successfully connected to the database.")

	downFlag := flag.Bool("down", false, "Run database migration down")
	downAllFlag := flag.Bool("down-all", false, "Run all database migrations down")

	// Parse the flags
	flag.Parse()

	if *downFlag {
		log.Println("Running database migration down...")
		migration.Down(db)
		log.Println("Successfully run database migration down.")
		return
	}

	if *downAllFlag {
		log.Println("Running all database migrations down...")
		migration.DownAll(db)
		log.Println("Successfully run all database migrations down.")
		return
	}

	log.Println("Running database migration up...")

	migration.Up(db)
	
	log.Println("Successfully run database migration up.")

	log.Println("Initializing validator...")
	
	SetupValidate()

	log.Println("Successfully initialized validator.")

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

	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusNotFound).JSON(map[string]any{
			"message": "Not Found",
			"data":    fmt.Sprintf("Not Found path %s", ctx.Path()),
		})
	})

	log.Println("Server started")
	log.Fatal(app.Listen(fmt.Sprintf(":%d", appPort)))
}