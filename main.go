package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/mojtabafarzaneh/hotel_reservation/api"
)

func main() {
	listenAdder := flag.String("listenAdder", ":5000", "the listen address of the api server")
	flag.Parse()
	app := fiber.New()
	apiv1 := app.Group("/api/v1")
	apiv1.Get("/user", api.HandleListUsers)
	apiv1.Get("/user/:id", api.HandelGetUser)
	app.Listen(*listenAdder)

}
