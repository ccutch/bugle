package bugle

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

// Audience of Members
type Audience struct {
	key *datastore.Key

	Name    string    `json:"name"`
	Owner   string    `json:"owner"`
	Created time.Time `json:"created"`
}

func (a Audience) KeyName() (name string) {
	name = strings.Replace(a.Name, " ", "_", -1)
	name = strings.ToLower(name)
	name = url.QueryEscape(name)
	return name
}

// Id gets the id of an audience
func (a Audience) ID() string {
	if i := a.key.ID; i != 0 {
		return strconv.Itoa(int(i))
	}

	return a.key.Name
}

// IsZero checks if all fields are zero values
func (a Audience) IsZero() bool { return a.Name == "" && a.Owner == "" && a.Created.IsZero() }

// Member of an Audience
type Member struct {
	key *datastore.Key
	aud *Audience

	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Joined time.Time `json:"joined"`
}
