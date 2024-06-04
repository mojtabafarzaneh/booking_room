package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mojtabafarzaneh/hotel_reservation/db"
	"github.com/mojtabafarzaneh/hotel_reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

type RoomHandler struct {
	store db.Store
}

func NewRoomHandler(store db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) BookingRoomHandler(c *fiber.Ctx) error {
	var prams BookingRoomParams
	if err := c.BodyParser(&prams); err != nil {
		return err
	}
	roomID := c.Params("id")
	roomOID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "internal server error",
		})
	}
	booking := types.Booking{
		RoomID:     roomOID,
		UserID:     user.ID,
		FromDate:   prams.FromDate,
		TillDate:   prams.TillDate,
		NumPersons: prams.NumPersons,
	}
	fmt.Println(booking)
	return nil
}
