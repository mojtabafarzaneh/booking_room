package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mojtabafarzaneh/booking_room/api"
	"github.com/mojtabafarzaneh/booking_room/db"
	"github.com/mojtabafarzaneh/booking_room/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var err error

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	mongoEndpoint := os.Getenv("MONGO_DB_URL_TEST")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.MongoDBNameEnvName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		Hotels:  hotelStore,
		Rooms:   db.NewMongoRoomStore(client, hotelStore),
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
	}
	user := fixtures.AddUser(store, "par", "jafar", false)
	fmt.Println("user token --->", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "moj", "far", true)
	fmt.Println("admin token --->", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "garndHotel", "rasht", 5, nil)
	room := fixtures.AddRoom(store, hotel.ID, 123.22, "small", false)
	booking := fixtures.AddBooking(store, room.ID, user.ID, 88, time.Now(), time.Now().AddDate(0, 0, 6))
	fmt.Println("booking ->", booking.ID)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("random hotel name %d", i)
		location := fmt.Sprintf("random hotel location %d", i)
		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}
}
