package example

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jakecoffman/gorest"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Author struct {
	ID    primitive.ObjectID `bson:"_id" json:"_id"`
	Name  string             `json:"name"`
	Books []Book             `json:"books"`
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
	*gorest.Controller
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
	id, _ := ctx.Get("id")
	var resource Author
	if err := ctx.BindJSON(&resource); err != nil {
		ctx.JSON(422, bson.M{"error": err.Error()})
		return
	}
	resource.ID = id.(primitive.ObjectID)
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
