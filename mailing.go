package bugle

import (
	"time"

	"cloud.google.com/go/datastore"
)

type Audience struct {
	key *datastore.Key

	Name    string    `json:"name"`
	Owner   string    `json:"owner"`
	Created time.Time `json:"created"`
}

func (a Audience) ID() string { return a.key.Name }

type Member struct {
	key *datastore.Key
	aud *Audience

	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Joined time.Time `json:"joined"`
}
