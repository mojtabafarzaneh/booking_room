package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mojtabafarzaneh/hotel_reservation/api"
	"github.com/mojtabafarzaneh/hotel_reservation/db"
	"github.com/mojtabafarzaneh/hotel_reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx          = context.Background()
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	client       *mongo.Client
)

func seedUser(isAdmin bool, fname, lname, email string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  fmt.Sprintf("%s_%s", fname, lname),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := userStore.InsertUsers(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s token ---> %s\n", user.Email, api.CreateTokenFromUser(user))
	return insertedUser

}

func seedHotel(name string, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)

	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel

}

func seedRoom(hotelID primitive.ObjectID, size string, price float64, seaside bool) *types.Room {
	room := types.Room{
		Size:    size,
		Price:   price,
		SeaSide: seaside,
		HotelID: hotelID,
	}
	insertedRoom, err := roomStore.InsertRoom(context.TODO(), &room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func seedBooking(uid, rid primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:   uid,
		RoomID:   rid,
		FromDate: from,
		TillDate: till,
	}
	insertedBooking, err := bookingStore.InsertBooking(context.TODO(), booking)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(booking.ID)
	return insertedBooking
}

func main() {
	seedUser(true, "moji", "farzaneh", "moj@gmail.com")
	user := seedUser(false, "ri", "stustu", "dk@gmail.com")
	seedHotel("grandhotel", "rasht", 3)
	seedHotel("hilton", "tehran", 4)
	hotel := seedHotel("abbasi", "isfahan", 5)
	room := seedRoom(hotel.ID, "small", 129.99, true)
	seedBooking(room.ID, user.ID, time.Now(), time.Now().AddDate(0, 0, 5))

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

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)

}
