package main

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"log"
	"github.com/gin-gonic/gin"
	"github.wwt.com/coffmanj/gorest"
	"github.wwt.com/coffmanj/gorest/example"
)

func main() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	db := session.DB("go-mongo-test")
	db.DropDatabase()

	bootstrap(db)

	r := router(db)
	r.Run("0.0.0.0:9889")
}

func bootstrap(db *mgo.Database) {
	books := []example.Book{{Title: "MyBook"}, {Title: "MyBook2"}}
	author := &example.Author{ID: bson.NewObjectId(), Name: "Bob", Books: books}

	err := db.C("author").Insert(author)
	if err != nil {
		log.Fatal(err)
	}

	authors := []example.Author{}
	err = db.C("author").Find(nil).All(&authors)
	if err != nil {
		log.Fatal(err)
	}
}

func router(db *mgo.Database) *gin.Engine {
	router := gin.Default()
	authorsRoute := router.Group("/authors")
	{
		ae := example.AuthorResource{MongoController: gorest.MongoController{
			C: db.C("author"), Resource: &example.Author{}}}
		authorsRoute.GET("/", ae.List)
		authorsRoute.GET("/:id", ae.Get)

		authorsRoute.POST("/", ae.Create)
		authorsRoute.PUT("/:id", ae.Update)
		authorsRoute.DELETE("/:id", ae.Delete)
	}
	return router
}
