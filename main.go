package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// Define a route for handling GET requests to "/employee".
	app.Get("/employee", func(c *fiber.Ctx) error {
		// Define an empty BSON document to query all documents in the "employees" collection.
		query := bson.D{{}}

		// Execute the query on the MongoDB server to retrieve all documents in the "employees" collection.
		cursor, err := mg.Db.Collection("employees").Find(c.Context(), query)
		if err != nil {
			// If an error occurs during the query, return a 500 Internal Server Error with the error message.
			return c.Status(500).SendString(err.Error())
		}

		// Create a slice to store the retrieved employee documents.
		var employees []Employee = make([]Employee, 0)

		// Decode all documents from the cursor and store them in the "employees" slice.
		if err := cursor.All(c.Context(), &employees); err != nil {
			// If an error occurs during decoding, return a 500 Internal Server Error with the error message.
			return c.Status(500).SendString(err.Error())
		}

		// Return the retrieved employee documents as a JSON response.
		return c.JSON(employees)
	})

	app.Post("/employee", func(c *fiber.Ctx) error {
		collection := mg.Db.Collection("employees")
		employee := new(Employee)
		if err := c.BodyParser(employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		employee.ID = ""
		insertionResult, err := collection.InsertOne(c.Context(), employee)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
		createdRecord := collection.FindOne(c.Context(), filter)
		createdEmployee := &Employee{}
		createdRecord.Decode(createdEmployee)
		return c.Status(201).JSON(createdEmployee)

	})
	app.Put("/employee/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		employeeID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		employee := new(Employee)
		if err := c.BodyParser(employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		query := bson.D{{Key: "_id", Value: employeeID}}
		update := bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{Key: "name", Value: employee.Name},
					{Key: "age", Value: employee.Age},
					{Key: "salary", Value: employee.Salary},
				},
			},
		}
		if err := mg.Db.Collection("employees").FindOneAndUpdate(c.Context(), query, update).Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				return c.SendStatus(400)
			}
			return c.SendStatus(500)
		}
		employee.ID = idParam
		return c.Status(200).JSON(employee)

	})

	app.Delete("/employee/id")
}
