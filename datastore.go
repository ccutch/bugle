package bugle

import (
	"context"

	"cloud.google.com/go/datastore"
)

func Datastore(id string) (*client, error) {
	c, err := datastore.NewClient(context.TODO(), id)
	return &client{c}, err
}

type client struct{ *datastore.Client }

func (c client) saveMailingList(list *MailingList) (path string, err error) {
	if list.key == nil {
		list.key = datastore.IncompleteKey("MailingList", nil)
	}

	list.key, err = c.Put(context.TODO(), list.key, list)
	return list.key.String(), err
}

func (c client) getMailingList(name string) (list MailingList, err error) {
	list.key = datastore.NameKey("MailingList", name, nil)
	err = c.Get(context.TODO(), list.key, &list)
	return list, err
}

func (c client) saveSubscription(sub *Subscription) (path string, err error) {
	sub.key = datastore.NameKey("Subscription", sub.Addr, sub.list.key)
	sub.key, err = c.Put(context.TODO(), sub.key, sub)
	return sub.key.String(), err
}

func (c client) getSubscriptions(list *MailingList) (subs []Subscription, err error) {
	subs = make([]Subscription, 0)
	q := datastore.NewQuery("Subscription").Ancestor(list.key)
	keys, err := c.GetAll(context.TODO(), q, &subs)
	for i, sub := range subs {
		sub.list = list
		sub.key = keys[i]
	}
	return subs, err
}

func (c client) deleteSubscription(list *MailingList, addr string) error {
	key := datastore.NameKey("Subscription", addr, list.key)
	return c.Delete(context.TODO(), key)
}
