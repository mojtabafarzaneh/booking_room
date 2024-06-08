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

type ResourceResp struct {
	Data   any `json:"data"`
	Result int `json:"result"`
	Page   int `json:"page"`
}

type HotelQueryParams struct {
	db.Pagination
	Rating int
}

func (h *HotelHandler) HandelGetHotels(c *fiber.Ctx) error {
	var params HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}

	filter := db.Map{
		"rating": params.Rating,
	}

	hotels, err := h.store.Hotels.GetHotels(c.Context(), filter, &params.Pagination)
	if err != nil {
		return err
	}

	resp := ResourceResp{
		Page:   int(params.Page),
		Result: int(params.Limit),
		Data:   hotels,
	}
	return c.JSON(resp)
}

func (h *HotelHandler) HandelGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}
	filter := bson.M{"hotelID": oid}

	rooms, err := h.store.Rooms.GetRoom(c.Context(), filter)
	if err != nil {
		return ErrInvalidID()
	}
	return c.JSON(rooms)

}

func (h *HotelHandler) HandelGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")

	hotel, err := h.store.Hotels.GetHOtelsByID(c.Context(), id)
	if err != nil {
		return ErrInvalidID()
	}
	return c.JSON(hotel)
}
