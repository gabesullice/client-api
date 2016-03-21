package crud

import (
	// stdlib
	"log"

	// external
	r "github.com/dancannon/gorethink"
)

var (
	session *r.Session
)

func init() {
	sess, err := r.Connect(r.ConnectOpts{
		Address:  "rethinkdb:28015",
		Database: "client_api",
		MaxIdle:  10,
		MaxOpen:  10,
	})
	if err != nil {
		log.Fatalln(err)
	}

	session = sess
}
