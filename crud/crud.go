package crud

import (
	// stdlib
	"log"
	"os"

	// external
	r "github.com/dancannon/gorethink"
)

var (
	session *r.Session
)

func init() {
	rethink_host := os.Getenv("RETHINKDB_HOST")

	sess, err := r.Connect(r.ConnectOpts{
		Address:  rethink_host + ":28015",
		Database: "client_api",
		MaxIdle:  10,
		MaxOpen:  10,
	})
	if err != nil {
		log.Fatalln(err)
	}

	session = sess
}
