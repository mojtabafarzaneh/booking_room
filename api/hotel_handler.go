package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mojtabafarzaneh/hotel_reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store db.Store
}

func NewHotelHandler(store db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandelGetHotels(c *fiber.Ctx) error {

	hotels, err := h.store.Hotels.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandelGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"hotelID": oid}

	rooms, err := h.store.Rooms.GetRoom(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)

}

func (h *HotelHandler) HandelGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")

	hotel, err := h.store.Hotels.GetHOtelsByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(hotel)
}
