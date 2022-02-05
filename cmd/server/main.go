package main

import (
	"net/http"

	"bugle.email"
)

func main() {
	http.ListenAndServe(":4000", bugle.NewRouter())
}
