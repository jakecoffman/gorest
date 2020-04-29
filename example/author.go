package example

import (
	"fmt"

	"github.com/jakecoffman/gorest"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Author struct {
	ID    primitive.ObjectID `bson:"_id" json:"_id"`
	Name  string             `json:"name"`
	Books []Book             `json:"books"`
}

func (a *Author) SetID(id primitive.ObjectID) {
	a.ID = id
}

func (a *Author) Decode(cursor gorest.Decoder) error {
	return cursor.Decode(a)
}

func (a *Author) Valid() error {
	if a.ID.IsZero() {
		return fmt.Errorf("author needs an `ID`")
	}
	if a.Name == "" {
		return fmt.Errorf("author needs a `Name`")
	}
	for _, book := range a.Books {
		if err := book.Valid(); err != nil {
			return err
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
