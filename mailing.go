package bugle

import "cloud.google.com/go/datastore"

type MailingList struct {
	key   *datastore.Key
	Name  string
	Owner string
}

type Subscription struct {
	key  *datastore.Key
	list *MailingList

	User string
	Addr string
}
