package crud

import (
	// stdlib
	"fmt"
	"log"
	"net/http"

	// internal
	"github.com/gabesullice/client-api/models"
	"github.com/gabesullice/client-api/storage"

	// external
	r "github.com/dancannon/gorethink"
	"github.com/manyminds/api2go"
)

type ContactResource struct {
	Session *r.Session
	Storage storage.Connection
}

func NewContactResource() ContactResource {
	return ContactResource{
		Session: session,
		Storage: storage.Connection{Session: session},
	}
}

func (c ContactResource) FindOne(ID string, req api2go.Request) (api2go.Responder, error) {
	var contact models.Contact
	if err := c.Storage.Get(ID, contact.GetName(), &contact); err != nil {
		return api2go.Response{}, api2go.NewHTTPError(err, "Unable to retrieve contact", http.StatusInternalServerError)
	}

	if contact.ID == "" {
		return api2go.Response{
			Code: http.StatusNotFound,
		}, nil
	}

	return api2go.Response{
		Res:  contact,
		Code: http.StatusOK,
	}, nil
}

func (c ContactResource) FindAll(req api2go.Request) (api2go.Responder, error) {
	var contacts []models.Contact
	if err := c.Storage.GetAll("contacts", &contacts); err != nil {
		return api2go.Response{}, api2go.NewHTTPError(err, "Unable to retrieve contacts", http.StatusInternalServerError)
	}

	if len(contacts) < 1 {
		return api2go.Response{
			Code: http.StatusNotFound,
		}, nil
	}

	return api2go.Response{
		Res:  contacts,
		Code: http.StatusOK,
	}, nil
}

func (c ContactResource) Create(obj interface{}, req api2go.Request) (api2go.Responder, error) {
	var resp api2go.Response

	contact, ok := obj.(models.Contact)
	if !ok {
		return resp, api2go.NewHTTPError(fmt.Errorf("Invalid instance given."), "Invalid instance given", http.StatusBadRequest)
	}

	res, err := r.Table("contacts").Insert(contact).RunWrite(c.Session)
	if err != nil {
		return resp, api2go.NewHTTPError(err, "Unable to save new resource", http.StatusInternalServerError)
	}

	contact.ID = res.GeneratedKeys[0]

	resp.Res = contact
	resp.Code = http.StatusCreated
	resp.Meta = map[string]interface{}{}

	return resp, nil
}

func (c ContactResource) Update(obj interface{}, req api2go.Request) (api2go.Responder, error) {
	var resp api2go.Response

	contact, ok := obj.(models.Contact)
	if !ok {
		return resp, api2go.NewHTTPError(fmt.Errorf("Invalid instance given."), "Invalid instance given", http.StatusBadRequest)
	}

	res, err := r.Table("contacts").Update(contact).RunWrite(c.Session)
	if err != nil {
		return resp, api2go.NewHTTPError(err, "Unable to update resource", http.StatusInternalServerError)
	}

	resp.Res = contact
	resp.Code = http.StatusNoContent
	resp.Meta = map[string]interface{}{
		"replaced": res.Replaced,
		"updated":  res.Updated,
	}

	return resp, nil
}

func (c ContactResource) Delete(ID string, req api2go.Request) (api2go.Responder, error) {
	var resp api2go.Response

	res, err := r.Table("contacts").Get(ID).Run(c.Session)
	if err != nil {
		return resp, api2go.NewHTTPError(err, "Could not find contact", http.StatusInternalServerError)
	}
	defer res.Close()

	if res.IsNil() {
		return resp, api2go.NewHTTPError(fmt.Errorf("Requested resource, %s, does not exist.", ID), "Resource not found.", http.StatusNotFound)
	}

	if _, err := r.Table("contacts").Get(ID).Delete().RunWrite(c.Session); err != nil {
		return resp, api2go.NewHTTPError(err, "Unable to delete resource", http.StatusInternalServerError)
	}

	resp.Code = http.StatusNoContent

	return resp, nil
}

func (c ContactResource) Search(search string) ([]models.Contact, error) {
	res, err := r.Table("contacts").Filter(func(user r.Term) r.Term {
		return user.Field("firstName").
			Match(fmt.Sprintf("(?i)%s", search))
	}).Run(c.Session)

	if err != nil {
		log.Printf("RethinkDB query failed. Error: %s", err.Error())
		return []models.Contact{}, err
	}

	var contacts []models.Contact
	if err := res.All(&contacts); err != nil {
		log.Printf("Unable to read contacts from result. Error: %s", err.Error())
		return []models.Contact{}, err
	}

	return contacts, nil
}
