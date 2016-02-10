package main

import (
	// stdlib
	//	"fmt"
	//	"log"
	"net/http"

	// internal
	"github.com/gabesullice/client-api/crud"
	"github.com/gabesullice/client-api/models"
	"github.com/gabesullice/client-api/storage"

	// external
	//"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go/jsonapi"
)

var (
	config = map[string]interface{}{
		"port":   "3000",
		"prefix": "v0",
	}

	resources = []resource{
		resource{Object: models.Contact{}, CRUD: crud.ContactResource{storage.NewDefaultStorage()}},
	}
)

type resource struct {
	Object jsonapi.MarshalIdentifier
	CRUD   api2go.CRUD
}

func main() {
	api := api2go.NewAPI(config["prefix"].(string))

	for _, resource := range resources {
		api.AddResource(resource.Object, resource.CRUD)
	}

	handler := api.Handler()

	http.ListenAndServe(":3000", handler)
}
