package main

import (
	"context"
	"flag"

	"github.com/adrianowsh/booking/db"
	"github.com/adrianowsh/booking/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi    = "mongodb://root:1q2w3e@localhost:27017"
	dbname   = "booking_db"
	userColl = "users"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	//handlers initialization
	userHandler := handlers.NewUserHandler(db.NewMongoUserStore(client, dbname))

	app := fiber.New(config)
	api := app.Group("/api")

	v1 := api.Group("/v1")

	v1.Get("/users", userHandler.HandleGetUsers)
	v1.Get("/users/:id", userHandler.HandleGetUser)
	v1.Put("/users/:id", userHandler.HandlePutUser)
	v1.Post("/users", userHandler.HandlerPostUser)
	v1.Delete("/users/:id", userHandler.HandleDeleteUser)

	app.Listen(*listenAddr)
}
