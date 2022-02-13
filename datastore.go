package bugle

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

func Datastore(id string) (*client, error) {
	c, err := datastore.NewClient(context.TODO(), id)
	return &client{c}, err
}

type client struct{ *datastore.Client }

func (c client) saveAudience(aud *Audience) (path string, err error) {
	if aud.key == nil {
		aud.key = datastore.IncompleteKey("Audience", nil)
	}

	aud.key, err = c.Put(context.TODO(), aud.key, aud)
	return aud.key.String(), err
}

func (c client) getAudience(name string) (aud Audience, err error) {
	aud.key = datastore.NameKey("Audience", name, nil)
	err = c.Get(context.TODO(), aud.key, &aud)
	return aud, err
}

func (c client) getAudienceForUser(u *user) (auds []Audience, err error) {
	q := datastore.NewQuery("Audience") //.Order("-Created").Filter("Owner =", u.Email)
	iter := c.Run(context.TODO(), q)

	for {
		var aud Audience
		key, ierr := iter.Next(&aud)
		aud.key = key
		if ierr == iterator.Done {
			break
		} else if ierr != nil {
			err = errors.Wrap(ierr, "Error fetching next audience")
			break
		}

		auds = append(auds, aud)
	}

	return auds, err
}

func (c client) saveMember(sub *Member) (path string, err error) {
	sub.key = datastore.NameKey("Member", sub.Email, sub.aud.key)
	sub.key, err = c.Put(context.TODO(), sub.key, sub)
	return sub.key.String(), err
}

func (c client) getMembers(aud *Audience) (members []Member, err error) {
	members = make([]Member, 0)
	q := datastore.NewQuery("Member").Ancestor(aud.key)
	keys, err := c.GetAll(context.TODO(), q, &members)
	for i, sub := range members {
		sub.aud = aud
		sub.key = keys[i]
	}
	return members, err
}

func (c client) deleteMember(aud *Audience, addr string) error {
	key := datastore.NameKey("Member", addr, aud.key)
	return c.Delete(context.TODO(), key)
}
