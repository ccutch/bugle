package main

import (
	"log"
	"net/http"

	"github.com/pkg/errors"

	"bugle.email"
)

func main() {
	log.Println("Server running @ http://localhost:4000")
	err := http.ListenAndServe(":4000", bugle.NewRouter())
	log.Fatal(errors.Wrap(err, "Failed to run server"))
}
