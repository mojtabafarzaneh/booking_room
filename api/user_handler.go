package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mojtabafarzaneh/hotel_reservation/types"
)

func HandleListUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "moji",
		LastName:  "farzaneh",
	}
	return c.JSON(u)
}

func HandelGetUser(c *fiber.Ctx) error {
	return c.JSON("01?")
}
