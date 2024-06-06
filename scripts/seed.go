package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mojtabafarzaneh/hotel_reservation/api"
	"github.com/mojtabafarzaneh/hotel_reservation/db"
	"github.com/mojtabafarzaneh/hotel_reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx    = context.Background()
	client *mongo.Client
)

func main() {
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
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().
		ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

}
