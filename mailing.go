package bugle

import "cloud.google.com/go/datastore"

type MailingList struct {
	key   *datastore.Key
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

type Subscription struct {
	key  *datastore.Key
	list *MailingList

	User string `json:"user"`
	Addr string `json:"address"`
}
