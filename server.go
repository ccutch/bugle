package bugle

import (
	"context"
	"net/http"
)

// Server creates a new router and binds all routes
func Server(client *client) http.Handler {
	s := &server{http.NewServeMux(), client}
	s.HandleFunc("/create", s.create)
	s.HandleFunc("/view", s.view)
	s.HandleFunc("/subs", s.subs)
	s.HandleFunc("/add", s.add)
	s.HandleFunc("/remove", s.remove)
	s.HandleFunc("/send", s.send)
	return s
}

// server composes http's ServeMux struct for http routing
type server struct {
	*http.ServeMux
	client *client
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, _ := s.client.getMailingList(r.URL.Query().Get("list"))
	ctx := context.WithValue(r.Context(), "list", list)
	s.ServeMux.ServeHTTP(w, r.WithContext(ctx))
}

// create creates a new list
// [POST] /create
func (s server) create(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.restrictMethods("POST")
	h.respond(db.saveMailingList(h.list()))
}

// view gets basic details for a list
// [GET] /view
func (s server) view(w http.ResponseWriter, r *http.Request) {
	h, _ := s.handle(w, r)
	h.restrictMethods("GET")
	h.respond(h.list())
}

// subs gets all subscriptsion for a list
// [GET] /subs
func (s server) subs(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.restrictMethods("GET")
	h.respond(db.getSubscriptions(h.list()))
}

// add creates a subscription for a list
// [POST] /add { listName, name, address }
func (s server) add(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.restrictMethods("POST")
	h.respond(db.saveSubscription(&Subscription{
		list: h.list(),
		User: r.FormValue("name"),
		Addr: r.FormValue("address"),
	}))
}

// remove removes a subscription from a list by address
// [DELETE] /add { listName, address }
func (s server) remove(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.restrictMethods("DELETE")
	h.respond(nil, db.deleteSubscription(
		h.list(),
		r.FormValue("address"),
	))
}

// send sends mail to all subscriptions on a list
// [POST] /send { message, listName }
func (s server) send(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.restrictMethods("POST")
	subs, err := db.getSubscriptions(h.list())
	h.respond(nil, err, h.gmail().SendEmail(h.body(), subs...))
}

// handle creates a new handler, we also defer a recovery func
// to handle serverErrors we encounter
func (s server) handle(w http.ResponseWriter, r *http.Request) (*handler, *client) {
	h := handler{w, r, nil, 200}
	defer func() {
		if r, ok := recover().(serverError); ok {
			h.code = r.code
			h.respond(nil, r.err)
		}
	}()
	return &h, s.client
}
