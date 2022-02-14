package bugle

import (
	"context"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

func Datastore(id string) (*database, error) {
	c, err := datastore.NewClient(context.Background(), id)
	return &database{c}, err
}

type database struct{ *datastore.Client }

func (db database) saveAudience(ctx context.Context, aud *Audience) (err error) {
	if ctx.Err() != nil {
		return
	}

	if aud.key == nil {
		aud.key = datastore.NameKey("Audience", aud.KeyName(), nil)
	}

	aud.key, err = db.Put(ctx, aud.key, aud)
	return err
}

func (db database) getAudience(ctx context.Context, name string) (aud Audience, err error) {
	if ctx.Err() != nil {
		return
	}

	if id, err := strconv.Atoi(name); err == nil {
		aud.key = datastore.IDKey("Audience", int64(id), nil)
	} else {
		aud.key = datastore.NameKey("Audience", name, nil)
	}

	err = db.Get(ctx, aud.key, &aud)

	return aud, err
}

func (db database) getAudienceForUser(ctx context.Context, u *user) (auds []Audience, err error) {
	if ctx.Err() != nil {
		return
	}

	q := datastore.NewQuery("Audience").Order("-Created").Filter("Owner =", u.Email)
	iter := db.Run(ctx, q)

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

func (db database) saveMember(ctx context.Context, sub *Member) (err error) {
	if ctx.Err() != nil {
		return
	}

	sub.key = datastore.NameKey("Member", sub.Email, sub.aud.key)
	sub.key, err = db.Put(ctx, sub.key, sub)
	return err
}

func (db database) getMembers(ctx context.Context, aud *Audience) (members []Member, err error) {
	if ctx.Err() != nil {
		return
	}

	members = make([]Member, 0)
	q := datastore.NewQuery("Member").Ancestor(aud.key)
	keys, err := db.GetAll(ctx, q, &members)
	for i, sub := range members {
		sub.aud = aud
		sub.key = keys[i]
	}
	return members, err
}

func (db database) deleteMember(ctx context.Context, aud *Audience, addr string) (err error) {
	if ctx.Err() != nil {
		return
	}

	key := datastore.NameKey("Member", addr, aud.key)
	return db.Delete(ctx, key)
}
