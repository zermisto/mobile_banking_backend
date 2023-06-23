package handlers


import (
	"github.com/gofiber/fiber/v2"
)

type Name2 struct {
	Name    string
	Message string
}

func UpdateUser(c *fiber.Ctx) error {

	Name := Name2{
		Name:    "Oat",
		Message: "success",
	}

	return c.JSON(Name)
}
