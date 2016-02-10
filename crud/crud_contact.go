package crud

import (
	// stdlib
	"net/http"

	// internal
	"github.com/gabesullice/client-api/models"
	"github.com/gabesullice/client-api/storage"

	// external
	"github.com/manyminds/api2go"
)

type ContactResource struct {
	ContactStorage storage.StorageInterface
}

func (c ContactResource) FindAll(req api2go.Request) (api2go.Responder, error) {
	var resp api2go.Response

	var contacts []models.Contact
	err := c.ContactStorage.FindAll("contacts", contacts)
	if err != nil {
		return resp, err
	}

	resp.Res = contacts
	resp.Code = http.StatusOK
	resp.Meta = map[string]interface{}{}

	return resp, nil
}

func (c ContactResource) FindOne(ID string, req api2go.Request) (api2go.Responder, error) {
	var resp api2go.Response

	contact, err := c.ContactStorage.FindOne("contacts", ID)
	if err != nil {
		return resp, err
	}

	resp.Res = contact
	resp.Code = http.StatusOK
	resp.Meta = map[string]interface{}{}

	return resp, nil
}

func (c ContactResource) Create(contact interface{}, req api2go.Request) (api2go.Responder, error) {
	var resp api2go.Response

	obj := contact.(models.Contact)
	obj, err := c.ContactStorage.Create(obj)
	if err != nil {
		return resp, err
	}

	resp.Res = obj
	resp.Code = http.StatusCreated
	resp.Meta = map[string]interface{}{}

	return resp, nil
}

func (c ContactResource) Delete(id string, req api2go.Request) (api2go.Responder, error) {
	var resp api2go.Response

	if err := c.ContactStorage.Delete("contacts", id); err != nil {
		return resp, err
	}

	resp.Res = id
	resp.Code = http.StatusNoContent
	resp.Meta = map[string]interface{}{}

	return resp, nil
}

func (c ContactResource) Update(contact interface{}, req api2go.Request) (api2go.Responder, error) {
	var resp api2go.Response

	err := c.ContactStorage.Update(contact.(models.Contact))
	if err != nil {
		return resp, err
	}

	resp.Res = contact
	resp.Code = http.StatusOK
	resp.Meta = map[string]interface{}{
		"replaced": 1,
	}

	return resp, nil
}
