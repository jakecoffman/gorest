package example

import (
	"github.com/jakecoffman/gorest"
	"gopkg.in/mgo.v2/bson"
)

type AuthorResource struct {
	gorest.MongoController
}

type Author struct {
	ID    bson.ObjectId `bson:"_id"`
	Name  string
	Books []Book
}

func (r *Author) New() gorest.Resource {
	return &Author{}
}

func (r *Author) NewList() interface{} {
	return &[]Author{}
}

func (a *Author) Id(id string) {
	a.ID = bson.ObjectIdHex(id)
}

func (a *Author) Valid() bool {
	isValid := string(a.ID) != "" && a.Name != ""
	if !isValid {
		return isValid
	}
	for _, book := range a.Books {
		if !book.Valid() {
			return false
		}
	}
	return true
}

type Book struct {
	Title string
}

func (b Book) Valid() bool {
	return b.Title != ""
}
