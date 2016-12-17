package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"gopkg.in/mgo.v2"
	"log"
	"github.com/gin-gonic/gin"
)

var db *mgo.Database

func init() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	db = session.DB("gorest-example")
	db.DropDatabase()

	gin.SetMode(gin.ReleaseMode)
}

func TestServer(t *testing.T) {
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/authors/", nil)
	if err != nil {
		t.Fatal(err)
	}

	router(db).ServeHTTP(resp, req)

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expected := `[]` + "\n"
	if string(bytes) != expected {
		t.Fatalf("%#v != %#v", string(bytes), expected)
	}
}
