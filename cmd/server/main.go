package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"

	"bugle.email"
)

func main() {
	db, err := bugle.Mongo(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(errors.Wrapf(err, "Failed to connect db: %s", os.Getenv("DB_URL")))
	}

	log.Println("Server running @ http://localhost:8080")
	err = http.ListenAndServe(":8080", bugle.Server(db))
	log.Fatal(errors.Wrap(err, "Failed to run server"))
}
