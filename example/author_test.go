package example_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jakecoffman/gorest/example"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

const (
	collection = "author"
)

func beforeEach() (*mongo.Database, *example.AuthorController) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("test-db")
	if err = db.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
	return db, &example.AuthorController{db.Collection(collection)}
}

func TestAuthorController_List(t *testing.T) {
	db, controller := beforeEach()
	actual := []interface{}{
		example.Author{ID: primitive.NewObjectID(), Name: "Alice"},
		example.Author{ID: primitive.NewObjectID(), Name: "Bob"},
	}
	if _, err := db.Collection(collection).InsertMany(context.Background(), actual); err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	controller.List(c)

	if data, err := json.Marshal(actual); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(data, resp.Body.Bytes()) {
		t.Error("Unexpected listing", string(resp.Body.Bytes()))
	}
}

func TestAuthorController_Create(t *testing.T) {
	db, controller := beforeEach()

	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	var err error
	c.Request, err = http.NewRequest("GET", "/", bytes.NewBufferString(`{"Name":"Test"}`))
	if err != nil {
		t.Fatal(err)
	}
	controller.Create(c)

	var actual example.Author
	if err = json.Unmarshal(resp.Body.Bytes(), &actual); err != nil {
		t.Fatal(err)
	}
	if actual.Name != "Test" || actual.ID.IsZero() || actual.Books != nil {
		t.Error("Unexpected listing", string(resp.Body.Bytes()))
	}
	cur, err := db.Collection(collection).Find(context.Background(), bson.M{})
	if err != nil {
		t.Fatal(err)
	}
	var count int
	for cur.Next(context.Background()) {
		count++
		var author example.Author
		if err = cur.Decode(&author); err != nil {
			t.Error(err)
		}
	}
	if count != 1 {
		t.Error("Not the right amount of results", count)
	}
}

func TestAuthorController_Get(t *testing.T) {
	db, controller := beforeEach()
	actual := example.Author{ID: primitive.NewObjectID(), Name: "Alice"}
	if _, err := db.Collection(collection).InsertOne(context.Background(), actual); err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{{Key: "id", Value: actual.ID.Hex()}}
	controller.Get(c)

	if data, err := json.Marshal(actual); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(data, resp.Body.Bytes()) {
		t.Error("Unexpected listing", string(resp.Body.Bytes()))
	}
}

func TestAuthorController_Update(t *testing.T) {
	db, controller := beforeEach()
	actual := example.Author{ID: primitive.NewObjectID(), Name: "Alice"}
	if _, err := db.Collection(collection).InsertOne(context.Background(), actual); err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	var err error
	c.Request, err = http.NewRequest("GET", "/", bytes.NewBufferString(`{"Name":"Bob"}`))
	if err != nil {
		t.Fatal(err)
	}
	c.Params = gin.Params{{Key: "id", Value: actual.ID.Hex()}}
	controller.Update(c)

	actual.Name = "Bob"
	if data, err := json.Marshal(actual); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(data, resp.Body.Bytes()) {
		t.Error("Unexpected listing", string(resp.Body.Bytes()))
	}
}

func TestAuthorController_Delete(t *testing.T) {
	db, controller := beforeEach()
	actual := example.Author{ID: primitive.NewObjectID(), Name: "Alice"}
	if _, err := db.Collection(collection).InsertOne(context.Background(), actual); err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{{Key: "id", Value: actual.ID.Hex()}}
	controller.Delete(c)

	cur, err := db.Collection(collection).Find(context.Background(), bson.M{})
	if err != nil {
		t.Fatal(err)
	}
	var count int
	for cur.Next(context.Background()) {
		count++
	}
	if count != 0 {
		t.Error("Not the right amount of results", count)
	}
}
