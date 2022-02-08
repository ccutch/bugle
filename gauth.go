package bugle

import "net/http"

func parseUser(*http.Request) *user {
	return &user{}
}

type user struct{}
