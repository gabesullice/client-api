package slack

import (
	// stdlib
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// internal
	"github.com/gabesullice/client-api/crud"
	"github.com/gabesullice/client-api/models"
)

type ResourceSearcher interface {
	Search(string) ([]models.Contact, error)
}

type SlackHandler struct {
	Resource ResourceSearcher
}

func NewSlackHandler() SlackHandler {
	return SlackHandler{
		Resource: crud.NewContactResource(),
	}
}

func (s SlackHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		log.Printf("Unable to parse form. Error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	search := req.FormValue("text")

	contacts, err := s.Resource.Search(search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(contacts) < 1 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No matches found for \"%s\"", search)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.Encode(newSlackMessage(contacts))
}
