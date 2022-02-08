package bugle

import (
	"context"
	"net/http"
)

// Server creates a new router and binds all routes
func Server(client *DBClient) http.Handler {
	s := &server{http.NewServeMux(), client}
	s.HandleFunc("/create", s.create)
	s.HandleFunc("/view", s.view)
	s.HandleFunc("/add", s.add)
	s.HandleFunc("/remove", s.remove)
	s.HandleFunc("/send", s.send)
	return s
}

// server composes http's ServeMux struct for http routing
type server struct {
	*http.ServeMux
	client *DBClient
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "list", r.URL.Query().Get("list"))
	s.ServeMux.ServeHTTP(w, r.WithContext(ctx))
}

// create creates a new list
// [POST] /create?list=<list>
func (s server) create(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.respond(db.NewList(h.listName()))
}

// view gets all subscriptsion for a list
// [GET] /view?list=<list>
func (s server) view(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.respond(db.GetSubscriptions(h.listName()))
}

// add creates a subscription for a list
// [POST] /add { listName, name, address }
func (s server) add(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.respond(db.NewSubscription(
		h.listName(),
		r.FormValue("name"),
		r.FormValue("address"),
	))
}

// remove removes a subscription from a list by address
// [POST] /add { listName, address }
func (s server) remove(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.respond(nil, db.DeleteSubscription(
		h.listName(),
		r.FormValue("address"),
	))
}

// send sends mail to all subscriptions on a list
// [POST] /send { message, listName }
func (s server) send(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	emails, err := db.GetSubscriptions(h.listName())
	h.respond(nil, err, h.gmail().SendEmail(h.body(), emails...))
}

// handle creates a new handler, we also defer a recovery func
// to handle serverErrors we encounter
func (s server) handle(w http.ResponseWriter, r *http.Request) (*handler, *DBClient) {
	defer func() {
		if r, ok := recover().(serverError); ok {
			http.Error(w, r.err.Error(), r.code)
		}
	}()

	return &handler{w, r}, s.client
}
