package db

import (
	"context"

	"github.com/mojtabafarzaneh/hotel_reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Map map[string]any
type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(ctx context.Context, filter Map, update Map) error
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection("hotels"),
	}
}

func (s *MongoHotelStore) Update(ctx context.Context, filter Map, update Map) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	resp, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = resp.InsertedID.(primitive.ObjectID)
	return hotel, nil
}
