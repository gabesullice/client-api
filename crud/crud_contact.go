package crud

import (
	// stdlib
	"fmt"
	"net/http"

	// internal
	"github.com/gabesullice/client-api/models"

	// external
	r "github.com/dancannon/gorethink"
	"github.com/manyminds/api2go"
)

type ContactResource struct {
	Session *r.Session
}

func NewContactResource() ContactResource {
	return ContactResource{Session: session}
}

func (c ContactResource) FindAll(req api2go.Request) (api2go.Responder, error) {
	var resp api2go.Response

	var contacts []models.Contact
	res, err := r.Table("contacts").Run(c.Session)
	if err != nil {
		return resp, api2go.NewHTTPError(err, "Could not find contacts", http.StatusInternalServerError)
	}
	defer res.Close()

	if err := res.All(&contacts); err != nil {
		return resp, api2go.NewHTTPError(err, "Could not read contacts", http.StatusInternalServerError)
	}

	resp.Res = contacts
	resp.Code = http.StatusOK
	resp.Meta = map[string]interface{}{}

	return resp, nil
}

func (c ContactResource) FindOne(ID string, req api2go.Request) (api2go.Responder, error) {
	var resp api2go.Response

	res, err := r.Table("contacts").Get(ID).Run(c.Session)
	if err != nil {
		return resp, api2go.NewHTTPError(err, "Could not find contact", http.StatusInternalServerError)
	}
	defer res.Close()

	if res.IsNil() {
		return resp, api2go.NewHTTPError(fmt.Errorf("Requested resource, %s, does not exist.", ID), "Resource not found.", http.StatusNotFound)
	}

	var contact models.Contact
	if err := res.One(&contact); err != nil {
		return resp, api2go.NewHTTPError(err, "Could not read contact", http.StatusInternalServerError)
	}

	resp.Res = contact
	resp.Code = http.StatusOK
	resp.Meta = map[string]interface{}{}

	return resp, nil
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
