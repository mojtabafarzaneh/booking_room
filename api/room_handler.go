package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mojtabafarzaneh/hotel_reservation/db"
	"github.com/mojtabafarzaneh/hotel_reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (p *BookingRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("connot book a room in the past")
	}
	return nil
}

type RoomHandler struct {
	store db.Store
}

func NewRoomHandler(store db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Rooms.GetRoom(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) BookingRoomHandler(c *fiber.Ctx) error {
	var prams BookingRoomParams
	if err := c.BodyParser(&prams); err != nil {
		return err
	}

	if err := prams.validate(); err != nil {
		return err
	}

	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "errors",
			Msg:  "internal server error",
		})
	}

	avilable, err := h.isRoomAvialable(c.Context(), roomID, prams)
	if err != nil {
		return err
	}
	if !avilable {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room %s is already booked", c.Params("id")),
		})
	}

	booking := types.Booking{
		RoomID:     roomID,
		UserID:     user.ID,
		FromDate:   prams.FromDate,
		TillDate:   prams.TillDate,
		NumPersons: prams.NumPersons,
	}

	insertedBooking, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	return c.JSON(insertedBooking)
}

func (h *RoomHandler) isRoomAvialable(ctx context.Context, roomID primitive.ObjectID, prams BookingRoomParams) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": prams.FromDate,
		},
		"tillDate": bson.M{
			"$lte": prams.TillDate,
		},
	}

	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil

}
