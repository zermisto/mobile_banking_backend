package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	Username string
	Password string
}

type Name struct {
	Name string
	User User
}

func CreateUser(c *fiber.Ctx) error {
	var req User
	c.BodyParser(&req)
	fmt.Println(req)
	fmt.Print("success")

	Name := Name{
		Name: "Oat",
		User: req,
	}

	return c.JSON(Name)
}
