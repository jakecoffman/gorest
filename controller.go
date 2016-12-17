package gorest

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

type Restful interface {
	List(ctx *gin.Context)
	Get(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type RestController struct {
	Restful
}

func (e *RestController) List(ctx *gin.Context) {
	ctx.JSON(415, bson.M{"error": "Method not allowed"})
}

func (e *RestController) Get(ctx *gin.Context) {
	ctx.JSON(415, bson.M{"error": "Method not allowed"})
}

func (e *RestController) Create(ctx *gin.Context) {
	ctx.JSON(415, bson.M{"error": "Method not allowed"})
}

func (e *RestController) Update(ctx *gin.Context) {
	ctx.JSON(415, bson.M{"error": "Method not allowed"})
}

func (e *RestController) Delete(ctx *gin.Context) {
	ctx.JSON(415, bson.M{"error": "Method not allowed"})
}
