package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mojtabafarzaneh/hotel_reservation/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("you have to be authenticated in order to access this page")
	}
	if !user.IsAdmin {
		return fmt.Errorf("you're not authorized to view this page")
	}
	return c.Next()
}
