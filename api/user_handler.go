package api

import (
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
		id = c.Params("id")
	)

	user, err := h.userStore.GetUserByID(c.Context(), id)

	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleListUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.Validate(); len(err) > 0 {
		return c.JSON(err)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUsers(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}
