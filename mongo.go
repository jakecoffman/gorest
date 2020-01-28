package gorest

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidIdFilter(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, bson.M{"error": "invalid bson ID"})
		return
	}
	ctx.Set("id", id)
}
