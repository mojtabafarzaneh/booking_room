package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/mojtabafarzaneh/booking_room/api"
	"github.com/mojtabafarzaneh/booking_room/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configuration
// 1. MongoDB endpoint
// 2. ListenAddress of our HTTP server
// 3. JWT secret
// 4. MongoDBName

var config = fiber.Config{
	// Override default error handler
	ErrorHandler: api.ErrorHandler,
}

func main() {

	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
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
	apiv1 := app.Group("/api/v1", api.JWTAuthentication(userStore))
	admin := apiv1.Group("/admin", api.AdminAuth)

	//auth handler
	auth.Post("/auth", authHandler.HandleAuthentication)

	//user handlers
	apiv1.Get("/user", userHandler.HandleListUsers)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user/:id", userHandler.HandelGetUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	//hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandelGetHotels)
	apiv1.Get("/hotel/:id/room", hotelHandler.HandelGetRooms)
	apiv1.Get("/hotel/:id/", hotelHandler.HandelGetHotel)

	//room handlers
	apiv1.Post("/rooms/:id/booking", roomHandler.BookingRoomHandler)
	apiv1.Get("/rooms", roomHandler.HandleGetRooms)

	//booking handlers
	admin.Get("/bookings", bookingHandler.HandelGetBookings)
	apiv1.Get("/bookings/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/bookings/:id/cancele", bookingHandler.HandleCancelBooking)

	listenAddress := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenAddress)

}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
