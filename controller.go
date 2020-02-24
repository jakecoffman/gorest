package gorest

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Controller struct {
	C       *mongo.Collection
	New     func() Resource
	Timeout time.Duration
	Limit   int64
}

func NewController(collection *mongo.Collection, resource interface{}) *Controller {
	return &Controller{
		C: collection,
		New: func() Resource {
			return reflect.New(reflect.TypeOf(resource)).Interface().(Resource)
		},
		Timeout: 60 * time.Second,
		Limit:   1000,
	}
}

type Resource interface {
	SetID(id primitive.ObjectID)
	Valid() error
	Decode(cursor *mongo.Cursor) error
}

func (r *Controller) List(ctx *gin.Context) {
	timeout, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

	cursor, err := r.C.Find(timeout, bson.M{}, options.Find().SetLimit(r.Limit))
	if err != nil {
		ctx.JSON(500, bson.M{"error": err.Error()})
		return
	}
	defer cursor.Close(timeout)

	results := reflect.New(reflect.SliceOf(reflect.TypeOf(r.New()))).Elem()
	for cursor.Next(timeout) {
		result := r.New()
		if err = result.Decode(cursor); err != nil {
			log.Println(err)
			ctx.JSON(500, bson.M{"error": "Decoding " + err.Error()})
			return
		}
		results = reflect.Append(results, reflect.ValueOf(result))
	}
	ctx.JSON(200, results.Interface())
}

func (r *Controller) Get(ctx *gin.Context) {
	timeout, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

	id, _ := ctx.Get("id")
	result := r.New()
	if cursor, err := r.C.Find(timeout, bson.M{"_id": id}); err != nil {
		ctx.JSON(404, bson.M{"error": err.Error()})
		return
	} else {
		defer cursor.Close(timeout)
		if err = result.Decode(cursor); err != nil {
			ctx.JSON(500, bson.M{"error": err.Error()})
			return
		}
	}
	ctx.JSON(200, result)
}

func (r *Controller) Delete(ctx *gin.Context) {
	timeout, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

	id, _ := ctx.Get("id")
	if result, err := r.C.DeleteOne(timeout, bson.M{"_id": id}); err != nil {
		ctx.JSON(500, bson.M{"error": err.Error()})
		return
	} else if result.DeletedCount == 0 {
		ctx.JSON(404, bson.M{"error": "not found"})
		return
	}
	ctx.JSON(204, nil)
}

func (r *Controller) Create(ctx *gin.Context) {
	timeout, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

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
	if _, err := r.C.InsertOne(timeout, resource); err != nil {
		ctx.JSON(500, bson.M{"error": err.Error()})
		return
	}
	ctx.JSON(201, resource)
}

func (r *Controller) Update(ctx *gin.Context) {
	timeout, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

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
	if _, err := r.C.ReplaceOne(timeout, bson.M{"_id": id}, resource); err != nil {
		ctx.JSON(404, bson.M{"error": err.Error()})
		return
	}
	ctx.JSON(200, resource)
}
