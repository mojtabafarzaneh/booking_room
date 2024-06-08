package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mojtabafarzaneh/booking_room/db"
	"github.com/mojtabafarzaneh/booking_room/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fname, lname string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fname, lname),
		Password:  fmt.Sprintf("%s_%s", fname, lname),
		FirstName: fname,
		LastName:  lname,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin

	insertedUser, err := store.User.InsertUsers(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser

}
func AddHotel(store *db.Store, name, location string, rating int, room []primitive.ObjectID) *types.Hotel {
	var roomIDS = room
	if room == nil {
		roomIDS = []primitive.ObjectID{}
	}
	hotel := &types.Hotel{
		Rooms:    roomIDS,
		Name:     name,
		Location: location,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotels.Insert(context.TODO(), hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddRoom(store *db.Store, hotelID primitive.ObjectID, price float64, size string, seaside bool) *types.Room {
	room := &types.Room{
		HotelID: hotelID,
		Price:   price,
		SeaSide: seaside,
		Size:    size,
	}
	insertedRoom, err := store.Rooms.InsertRoom(context.TODO(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddBooking(store *db.Store, roomID, userID primitive.ObjectID, numPersons int, fDate, tDate time.Time) *types.Booking {
	booking := &types.Booking{
		RoomID:     roomID,
		UserID:     userID,
		NumPersons: numPersons,
		FromDate:   fDate,
		TillDate:   tDate,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.TODO(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}
