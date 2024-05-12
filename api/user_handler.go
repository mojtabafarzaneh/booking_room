package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/mojtabafarzaneh/hotel_reservation/db"
	"github.com/mojtabafarzaneh/hotel_reservation/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(usertore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: usertore,
	}
}

func (h *UserHandler) HandelGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = context.Background()
	)

	user, err := h.userStore.GetUserByID(ctx, id)

	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandelListUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "moji",
		LastName:  "farzaneh",
	}
	return c.JSON(u)
}
