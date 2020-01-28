package gorest

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Controller struct {
	C *mongo.Collection
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
