package api

import (
	"github.com/germanx/hotel-reservation/types"

	"github.com/gofiber/fiber/v2"
)

func HandleGetUser(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "James",
		LastName:  "Watercooler",
	}
	return c.JSON(u)
}

func HandleGetUsers(c *fiber.Ctx) error {
	return c.JSON("James")
}
