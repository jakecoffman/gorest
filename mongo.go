package gorest

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

type MongoController struct {
	RestController
	C *mgo.Collection
	Resource Resource
}

func (r *MongoController) List(ctx *gin.Context) {
	a := r.Resource.NewList()
	if err := r.C.Find(nil).All(a); err != nil {
		ctx.JSON(500, bson.M{"error": err})
		return
	}
	ctx.JSON(200, a)
}

func (r *MongoController) Get(ctx *gin.Context) {
	id := bson.ObjectIdHex(ctx.Param("id"))
	a := r.Resource.New()
	if err := r.C.FindId(id).One(a); err != nil {
		ctx.JSON(404, bson.M{"error": err.Error()})
		return
	}
	ctx.JSON(200, a)
}

func (r *MongoController) Create(ctx *gin.Context) {
	resource := r.Resource.New()
	if err := ctx.BindJSON(resource); err != nil {
		ctx.JSON(422, bson.M{"error": err.Error()})
		return
	}
	resource.Id(bson.NewObjectId().Hex())
	if !resource.Valid() {
		ctx.JSON(400, bson.M{"error": "invalid resource"})
		return
	}
	if err := r.C.Insert(resource); err != nil {
		ctx.JSON(500, bson.M{"error": err.Error()})
		return
	}
	ctx.JSON(201, resource)
}

func (r *MongoController) Update(ctx *gin.Context) {
	id := bson.ObjectIdHex(ctx.Param("id"))
	resource := r.Resource.New()
	if err := ctx.BindJSON(resource); err != nil {
		ctx.JSON(422, bson.M{"error": err.Error()})
		return
	}
	resource.Id(id.Hex())
	if !resource.Valid() {
		ctx.JSON(400, bson.M{"error": "invalid resource"})
		return
	}
	if err := r.C.UpdateId(id, resource); err != nil {
		ctx.JSON(404, bson.M{"error": "not found"})
		return
	}
	ctx.JSON(200, resource)
}

func (r *MongoController) Delete(ctx *gin.Context) {
	id := bson.ObjectIdHex(ctx.Param("id"))
	if err := r.C.RemoveId(id); err != nil {
		ctx.JSON(404, bson.M{"error": err.Error()})
		return
	}
	ctx.JSON(200, bson.M{})
}
