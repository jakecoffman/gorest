package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jakecoffman/gorest/example"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("go-mongo-test")

	bootstrap(db)

	router := gin.Default()
	authorsRoute := router.Group("/authors")
	{
		ae := example.AuthorController {
			C: db.Collection("author"),
		}
		authorsRoute.GET("", ae.List)
		authorsRoute.GET("/:id", ae.Get)

		authorsRoute.POST("", ae.Create)
		authorsRoute.PUT("/:id", ae.Update)
		authorsRoute.DELETE("/:id", ae.Delete)
	}
	if err = router.Run("localhost:9889"); err != nil {
		log.Fatal(err)
	}
}

func bootstrap(db *mongo.Database) {
	books := []example.Book{{Title: "MyBook"}, {Title: "MyBook2"}}
	author := &example.Author{ID: primitive.NewObjectID(), Name: "Bob", Books: books}

	_, err := db.Collection("author").InsertOne(context.Background(), author)
	if err != nil {
		log.Fatal(err)
	}

	authors := []example.Author{}
	cursor, err := db.Collection("author").Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	for cursor.Next(context.Background()) {
		var author example.Author
		if err = cursor.Decode(&author); err != nil {
			log.Fatal(err)
		}
		authors = append(authors, author)
	}
}
