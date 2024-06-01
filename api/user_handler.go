package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/mojtabafarzaneh/hotel_reservation/db"
	"github.com/mojtabafarzaneh/hotel_reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	store db.Store
}

func NewUserHandler(userStore db.Store) *UserHandler {
	return &UserHandler{
		store: userStore,
	}
}

func (h *UserHandler) HandelGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)

	user, err := h.store.User.GetUserByID(c.Context(), id)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "not found"})
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleListUsers(c *fiber.Ctx) error {
	users, err := h.store.User.GetUsers(c.Context())
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
	insertedUser, err := h.store.User.InsertUsers(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		params types.UpdateUserParams
		userID = c.Params("id")
	)
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	if err := c.BodyParser(&params); err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	if err := h.store.User.USerUpdate(c.Context(), filter, params); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": userID})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := h.store.User.DeleteUser(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": userID})
}
