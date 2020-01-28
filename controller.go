package gorest

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Controller struct {
	C *mongo.Collection
	New func() Resource
}

type Resource interface {
	SetID(id primitive.ObjectID)
	Valid() error
}

func (r *Controller) List(ctx *gin.Context) {
	cursor, err := r.C.Find(context.Background(), bson.M{})
	if err != nil {
		ctx.JSON(500, bson.M{"error": err})
		return
	}
	results := []map[string]interface{}{}
	for cursor.Next(context.Background()) {
		var result map[string]interface{}
		if err = cursor.Decode(&result); err != nil {
			log.Println(err)
			ctx.JSON(500, bson.M{"error": "Decoding"})
			return
		}
		results = append(results, result)
	}
	ctx.JSON(200, results)
}

func (r *Controller) Get(ctx *gin.Context) {
	id, _ := ctx.Get("id")
	var result map[string]interface{}
	if err := r.C.FindOne(context.Background(), bson.M{"_id": id}).Decode(&result); err != nil {
		ctx.JSON(404, bson.M{"error": err.Error()})
		return
	}
	ctx.JSON(200, result)
}

func (r *Controller) Delete(ctx *gin.Context) {
	id, _ := ctx.Get("id")
	if _, err := r.C.DeleteOne(context.Background(), bson.M{"_id": id}); err != nil {
		ctx.JSON(404, bson.M{"error": err.Error()})
		return
	}
	ctx.JSON(204, nil)
}

func (r *Controller) Create(ctx *gin.Context) {
	resource := r.New()
	if err := ctx.BindJSON(&resource); err != nil {
		ctx.JSON(422, bson.M{"error": err.Error()})
		return
	}
	resource.SetID(primitive.NewObjectID())
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

func (r *Controller) Update(ctx *gin.Context) {
	id, _ := ctx.Get("id")
	resource := r.New()
	if err := ctx.BindJSON(&resource); err != nil {
		ctx.JSON(422, bson.M{"error": err.Error()})
		return
	}
	resource.SetID(id.(primitive.ObjectID))
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
