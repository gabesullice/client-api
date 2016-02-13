package models

import (
	// stdlib
	"errors"
	"fmt"
	//"log"

	// external
	"github.com/manyminds/api2go/jsonapi"
)

type Phone struct {
	Type      string `gorethink:"type" json:"type"`
	Number    string `gorethink:"number" json:"number"`
	Extension string `gorethink:"extension" json:"extension,omitempty"`
}

type Contact struct {
	ID        string   `gorethink:"id,omitempty" json:"-"`
	Name      string   `gorethink:"name" json:"name"`
	FirstName string   `gorethink:"firstName" json:"firstName"`
	LastName  string   `gorethink:"lastName" json:"lastName"`
	Synonyms  []string `gorethink:"synonyms" json:"synonyms"`
	Position  string   `gorethink:"position" json:"position"`
	Phones    []Phone  `gorethink:"phones" json:"phones"`
	Emails    []string `gorethink:"emails" json:"emails"`
	Notes     string   `gorethink:"notes" json:"notes"`
	Related   []string `gorethink:"related" json:"related"`
}

func (c Contact) GetName() string {
	return "contacts"
}

func (c *Contact) SetID(ID string) error {
	c.ID = ID
	return nil
}

func (c Contact) GetID() string {
	return c.ID
}

func (c Contact) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: c.GetName(),
			Name: "related",
		},
	}
}

func (c Contact) GetReferencedIDs() []jsonapi.ReferenceID {
	var refs []jsonapi.ReferenceID
	for k := range c.Related {
		refs = append(refs, jsonapi.ReferenceID{
			ID:   c.Related[k],
			Type: c.GetName(),
			Name: "related",
		})
	}
	return refs
}

func (c *Contact) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "related" {
		c.Related = IDs
		return nil
	}

	return errors.New(fmt.Sprintf("There is no %s relationship on %s", name, c.GetName()))
}

//func (c *Contact) AddToManyIDs(name string, IDs []string) error {
//	if name == "related" {
//		c.Related = append(c.Related, IDs...)
//		return nil
//	}
//
//	return errors.New(fmt.Sprintf("There is no %s relationship on %s", name, c.GetName()))
//}
//
//func (c *Contact) DeleteToManyIDs(name string, IDs []string) error {
//	if name == "related" {
//		for _, ID := range IDs {
//			for pos, oldID := range c.Related {
//				if ID == oldID {
//					// match, this ID must be removed
//					c.Related = append(c.Related[:pos], c.Related[pos+1:]...)
//				}
//			}
//		}
//	}
//
//	return errors.New(fmt.Sprintf("There is no %s relationship on %s", name, c.GetName()))
//}
