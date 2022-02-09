package main

import (
	"log"
	"net/http"

	"github.com/pkg/errors"

	"bugle.email"
)

func main() {
	db, err := bugle.Datastore("bugle-340607")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to connect db"))
	}

	log.Println("Server running @ http://localhost:8080")
	err = http.ListenAndServe(":8080", bugle.Server(db))
	log.Fatal(errors.Wrap(err, "Failed to run server"))
}
