package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mojtabafarzaneh/hotel_reservation/api"
	"github.com/mojtabafarzaneh/hotel_reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://192.168.1.161:27017"
const dbname = "hotel_reservation"
const userColl = "users"

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	coll := client.Database(dbname).Collection(userColl)

	user := types.User{
		FirstName: "moji",
		LastName:  "farzaneh",
	}
	res, err := coll.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

	listenAdder := flag.String("listenAdder", ":5000", "the listen address of the api server")
	flag.Parse()
	app := fiber.New()
	apiv1 := app.Group("/api/v1")
	apiv1.Get("/user", api.HandleListUsers)
	apiv1.Get("/user/:id", api.HandelGetUser)
	app.Listen(*listenAdder)

}
