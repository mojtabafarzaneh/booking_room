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

func seedHotel(name, location string) error {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
	}

	room := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 99.9,
		},
		{
			Type:      types.DeluxeRoomtype,
			BasePrice: 299.9,
		},
		{
			Type:      types.SeasideRoomtype,
			BasePrice: 199.9,
		},
		{
			Type:      types.DoubleRoomType,
			BasePrice: 122.9,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)

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
	seedHotel("grandhotel", "rasht")
	seedHotel("hilton", "tehran")
	seedHotel("abbasi", "isfahan")
	seedHotel("the grand budapest hotel", "budapest")

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
