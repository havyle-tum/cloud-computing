package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Defines a "model" that we can use to communicate with the
// frontend or the database
type Book struct {
	MongoID     primitive.ObjectID `bson:"_id,omitempty"`
	ID          string
	BookName    string
	BookAuthor  string
	BookEdition string
	BookPages   string
	BookYear    string
}

func deleteBook(coll *mongo.Collection, id string) error {
	filter := bson.M{"id": id}
	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("book with ID %s not found", id)
	}

	return nil
}

// Here we make sure the connection to the database is correct and initial
// configurations exists. Otherwise, we create the proper database and collection
// we will store the data.
// To ensure correct management of the collection, we create a return a
// reference to the collection to always be used. Make sure if you create other
// files, that you pass the proper value to ensure communication with the
// database
func prepareDatabase(client *mongo.Client, dbName string, collecName string) (*mongo.Collection, error) {
	db := client.Database(dbName)

	names, err := db.ListCollectionNames(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	if !slices.Contains(names, collecName) {
		cmd := bson.D{{Key: "create", Value: collecName}}
		var result bson.M
		if err = db.RunCommand(context.TODO(), cmd).Decode(&result); err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	coll := db.Collection(collecName)
	return coll, nil
}

func main() {
	// Connect to the database. Such defer keywords are used once the local
	// context returns; for this case, the local context is the main function
	// By user defer function, we make sure we don't leave connections
	// dangling despite the program crashing. Isn't this nice? :D
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// TODO: make sure to pass the proper username, password, and port
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	mongoURI := os.Getenv("DATABASE_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))

	// This is another way to specify the call of a function. You can define inline
	// functions (or anonymous functions, similar to the behavior in Python)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// You can use such name for the database and collection, or come up with
	// one by yourself!
	coll, err := prepareDatabase(client, "exercise-1", "information")

	// Here we prepare the server
	e := echo.New()

	// Log the requests. Please have a look at echo's documentation on more
	// middleware
	e.Use(middleware.Logger())

	e.DELETE("/api/books/:id", func(c echo.Context) error {
		id := c.Param("id")

		if err := deleteBook(coll, id); err != nil {
			if err.Error() == fmt.Sprintf("book with ID %s not found", id) {
				return c.JSON(http.StatusNotFound, map[string]string{
					"error": err.Error(),
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "Book deleted successfully",
		})
	})

	// We start the server and bind it to port 3030. For future references, this
	// is the application's port and not the external one. For this first exercise,
	// they could be the same if you use a Cloud Provider. If you use ngrok or similar,
	// they might differ.
	// In the submission website for this exercise, you will have to provide the internet-reachable
	// endpoint: http://<host>:<external-port>
	fmt.Println("Books DELETE service starting on port 8084")
	e.Logger.Fatal(e.Start(":8084"))
}