package main

import (
	"context"
	"log"

	"github.com/mojtabafarzaneh/hotel_reservation/db"
	"github.com/mojtabafarzaneh/hotel_reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx        = context.Background()
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	client     *mongo.Client
)

func seedHotel(name string, location string, rating int) error {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	room := []types.Room{
		{
			Size:  "small",
			Price: 99.9,
		},
		{
			Size:  "kingSize",
			Price: 299.9,
		},
		{
			Size:  "normal",
			Price: 199.9,
		},
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)

	if err != nil {
		log.Fatal(err)
	}
	for _, rooms := range room {
		rooms.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &rooms)
		if err != nil {
			log.Fatal(err)
		}

	}
	return nil

}

func main() {
	seedHotel("grandhotel", "rasht", 3)
	seedHotel("hilton", "tehran", 4)
	seedHotel("abbasi", "isfahan", 5)

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

}
