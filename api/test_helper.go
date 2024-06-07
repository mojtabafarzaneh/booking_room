package api

import (
	"context"
	"log"
	"testing"

	"github.com/mojtabafarzaneh/hotel_reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	*db.Store
	client *mongo.Client
}

func (tdb *testdb) teardrop(t *testing.T) {
	if err := tdb.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}

}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(db.DBURI))
	if err != nil {
		t.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	return &testdb{
		client: client,
		Store: &db.Store{
			Hotels:  hotelStore,
			Rooms:   db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client),
			User:    db.NewMongoUserStore(client),
		},
	}
}
