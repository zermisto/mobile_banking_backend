package main

import (
	"boilerplate/database"
	"boilerplate/handlers"

	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	//"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	port = flag.String("port", ":3000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {
	fmt.Printf("Secret username is: %s\n", os.Getenv("USERNAME"))
	fmt.Printf("Secret password is: %s\n", os.Getenv("PASSWORD"))
	// Parse command-line flags
	flag.Parse()

	// Connected with database
	database.Connect()

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: *prod, // go run app.go -prod
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		TimeZone: "Asia/Bangkok",
	}))

	app.Post("/user", handlers.CreateUser)
	app.Get("/users", handlers.UpdateUser)

	app.Get("/parent/:id", handlers.Parent{}.Find)
	app.Get("/parent", handlers.Parent{}.FindAll)
	app.Post("/parent", handlers.CreateParentHandler)

	app.Get("/student", handlers.Student{}.FindAll)
	app.Post("/student", handlers.CreateStudentHandler)

	app.Get("/payment", handlers.Payment{}.Find)
	app.Post("/payment", handlers.CreatePaymentHandler)
	app.Put("/payment/:id", handlers.PaidHandler)

	/*// Create a /api/v1 endpoint
	v1 := app.Group("/api/v1")

	// Bind handlers
	v1.Get("/users", handlers.UserList)
	v1.Post("/users", handlers.UserCreate)*/

	// app.Use(basicauth.New(basicauth.Config{
	// 	Users: map[string]string{
	// 		os.Getenv("USERNAME"): os.Getenv("PASSWORD"),
	// 	},
	// }))
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	// Setup static files
	app.Static("/", "./static/public")

	// Handle not founds
	app.Use(handlers.NotFound)

	// Listen on port 3000
	log.Fatal(app.Listen(*port)) // go run app.go -port=:3000
}
