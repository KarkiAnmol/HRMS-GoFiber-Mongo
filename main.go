package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg MongoInstance

const dbName = "fiber-hrms"
const mongoURI = "mongodb://localhost:27017" + dbName

type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    int     `json:"age"`
}

// connectDB establishes a connection to the MongoDB server.
// It returns an error if the connection cannot be established.
func connectDB() error {
	// Set a timeout for the database connection attempt.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	// Connect to the MongoDB server using the provided URI.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		// If an error occurs during connection, return the error.
		return err
	}
	// Ensure the context is canceled when the function exits.
	defer cancel()

	// Select the specific database from the connected MongoDB server.
	db := client.Database(dbName)

	// Store the MongoDB client and selected database in a global variable.
	// This allows other parts of the program to access the same MongoDB instance.
	mg = MongoInstance{
		Client: client,
		Db:     db,
	}

	// Connection successful, return nil (no error).
	return nil
}

func main() {
	if err := connectDB(); err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	app.Get("/employee", func(c *fiber.Ctx) error {
		query := bson.D{{}}
		cursor, err := mg.Db.Collection("employees").Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		var employees []Employee = make([]Employee, 0)
		if err := cursor.All(c.Context(), &employees); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(employees)

	})
	app.Post("/employee")
	app.Put("/employee/:id")

	app.Delete("/employee/id")
}
