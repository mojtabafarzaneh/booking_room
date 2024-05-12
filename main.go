package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mojtabafarzaneh/hotel_reservation/api"
	"github.com/mojtabafarzaneh/hotel_reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://192.168.1.161:27017"

var config = fiber.Config{
	// Override default error handler
	ErrorHandler: func(c *fiber.Ctx, err error) error {

		return c.JSON(map[string]string{"err": err.Error()})
	},
}

func main() {

	listenAdder := flag.String("listenAdder", ":5000", "the listen address of the api server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	//handler initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")
	apiv1.Get("/user", userHandler.HandleListUsers)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user/:id", userHandler.HandelGetUser)
	app.Listen(*listenAdder)

}
