package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mojtabafarzaneh/hotel_reservation/api"
	"github.com/mojtabafarzaneh/hotel_reservation/api/middleware"
	"github.com/mojtabafarzaneh/hotel_reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	// Override default error handler
	ErrorHandler: func(c *fiber.Ctx, err error) error {

		return c.JSON(map[string]string{"err": err.Error()})
	},
}

func main() {

	listenAdder := flag.String("listenAdder", ":5000", "the listen address of the api server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	//handler initialization
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	userStore := db.NewMongoUserStore(client)
	bookingStore := db.NewMongoBookingStore(client)
	store := &db.Store{
		Hotels:  hotelStore,
		Rooms:   roomStore,
		User:    userStore,
		Booking: bookingStore,
	}
	hotelHandler := api.NewHotelHandler(*store)
	userHandler := api.NewUserHandler(*store)
	authHandler := api.NewAuthHnadler(userStore)
	roomHandler := api.NewRoomHandler(*store)
	bookingHandler := api.NewBookingHandler(*store)

	app := fiber.New(config)
	auth := app.Group("/api")
	apiv1 := app.Group("/api/v1", middleware.JWTAuthentication(userStore))
	admin := apiv1.Group("/admin", middleware.AdminAuth)

	//auth handler
	auth.Post("/auth", authHandler.HandleAuthentication)

	//user handlers
	apiv1.Get("/user", userHandler.HandleListUsers)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user/:id", userHandler.HandelGetUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	//hotel handlers
	apiv1.Get("/hotels", hotelHandler.HandelGetHotels)
	apiv1.Get("/hotel/:id/room", hotelHandler.HandelGetRooms)
	apiv1.Get("/hotel/:id/", hotelHandler.HandelGetHotel)

	//room handlers
	apiv1.Post("/rooms/:id/booking", roomHandler.BookingRoomHandler)
	apiv1.Get("/rooms", roomHandler.HandleGetRooms)

	//booking handlers
	admin.Get("/bookings", bookingHandler.HandelGetBookings)
	apiv1.Get("/bookings/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/bookings/:id/cancele", bookingHandler.HandleCancelBooking)

	app.Listen(*listenAdder)

}
