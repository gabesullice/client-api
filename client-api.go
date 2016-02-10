package main

import (
	// stdlib
	//	"fmt"
	"log"
	"net/http"

	// internal
	"github.com/gabesullice/client-api/crud"
	"github.com/gabesullice/client-api/models"

	// external
	//"github.com/julienschmidt/httprouter"
	r "github.com/dancannon/gorethink"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go/jsonapi"
)

var (
	config = map[string]interface{}{
		"port":   "3000",
		"prefix": "v0",
	}

	session *r.Session

	resources = []resource{
		resource{Object: models.Contact{}, CRUD: crud.NewContactResource()},
	}
)

type resource struct {
	Object jsonapi.MarshalIdentifier
	CRUD   api2go.CRUD
}

func init() {
	sess, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "client-api",
		MaxIdle:  10,
		MaxOpen:  10,
	})
	if err != nil {
		log.Fatalln(err)
	}

	session = sess
}

func main() {
	api := api2go.NewAPI(config["prefix"].(string))

	for _, resource := range resources {
		api.AddResource(resource.Object, resource.CRUD)
	}

	handler := api.Handler()

	log.Printf("Awaiting connections on %s...", ":3000")
	http.ListenAndServe(":3000", handler)
}
