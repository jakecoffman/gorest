package example

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Author struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string
	Books []Book
}

func (a *Author) Valid() error {
	if a.ID.IsZero() {
		return fmt.Errorf("author needs an `ID`")
	}
	if a.Name == "" {
		return fmt.Errorf("author needs a `Name`")
	}
	for _, book := range a.Books {
		if book.Valid() != nil {
			return book.Valid()
		}
	}
	return nil
}

type Book struct {
	Title string
}

func (b Book) Valid() error {
	if b.Title == "" {
		return fmt.Errorf("book needs a title")
	}
	return nil
}

type AuthorController struct {
	C *mongo.Collection
}

func (r *AuthorController) List(ctx *gin.Context) {
	cursor, err := r.C.Find(context.Background(), bson.M{})
	if err != nil {
		ctx.JSON(500, bson.M{"error": err})
		return
	}
	results := []Author{}
	for cursor.Next(context.Background()) {
		var result Author
		if err = cursor.Decode(&result); err != nil {
			log.Println(err)
			ctx.JSON(500, bson.M{"error": "Decoding"})
			return
		}
		results = append(results, result)
	}
	ctx.JSON(200, results)
}

func (r *AuthorController) Get(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, bson.M{"error": "invalid bson ID"})
		return
	}
	var result Author
	if err := r.C.FindOne(context.Background(), bson.M{"_id": id}).Decode(&result); err != nil {
		ctx.JSON(404, bson.M{"error": err.Error()})
		return
	}
	ctx.JSON(200, result)
}

func (r *AuthorController) Create(ctx *gin.Context) {
	var resource Author
	if err := ctx.BindJSON(&resource); err != nil {
		ctx.JSON(422, bson.M{"error": err.Error()})
		return
	}
	resource.ID = primitive.NewObjectID()
	if resource.Valid() != nil {
		ctx.JSON(400, bson.M{"error": "invalid resource: " + resource.Valid().Error()})
		return
	}
	if _, err := r.C.InsertOne(context.Background(), resource); err != nil {
		ctx.JSON(500, bson.M{"error": err.Error()})
		return
	}
	ctx.JSON(201, resource)
}

func (r *AuthorController) Update(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, bson.M{"error": "invalid bson ID"})
		return
	}
	var resource Author
	if err := ctx.BindJSON(&resource); err != nil {
		ctx.JSON(422, bson.M{"error": err.Error()})
		return
	}
	resource.ID = id
	if resource.Valid() != nil {
		ctx.JSON(400, bson.M{"error": "invalid resource: " + resource.Valid().Error()})
		return
	}
	if _, err := r.C.ReplaceOne(context.Background(), bson.M{"_id": id}, resource); err != nil {
		ctx.JSON(404, bson.M{"error": err.Error()})
		return
	}
	ctx.JSON(200, resource)
}

func (r *AuthorController) Delete(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, bson.M{"error": "invalid bson id"})
		return
	}
	if _, err := r.C.DeleteOne(context.Background(), bson.M{"_id": id}); err != nil {
		ctx.JSON(404, bson.M{"error": err.Error()})
		return
	}
	ctx.JSON(204, nil)
}
