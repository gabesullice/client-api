package storage

import (
	// external
	r "github.com/dancannon/gorethink"
)

type Connection struct {
	Session *r.Session
}

type ObjectNamer interface {
	GetName() string
}

func (c Connection) Get(ID string, name string, obj interface{}) error {
	res, err := r.Table(name).Get(ID).Run(c.Session)
	if err != nil {
		return err
	}
	defer res.Close()

	if res.IsNil() {
		return nil
	}

	if err := res.One(obj); err != nil {
		return err
	}

	return nil
}

func (c Connection) GetAll(name string, objects interface{}) error {
	res, err := r.Table(name).Run(c.Session)
	if err != nil {
		return err
	}
	defer res.Close()

	if res.IsNil() {
		return nil
	}

	if err := res.All(objects); err != nil {
		return err
	}

	return nil
}
