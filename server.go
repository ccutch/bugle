package bugle

import (
	"context"
	"net/http"
)

// Server creates a new router and binds all routes
func Server(client *client) http.Handler {
	s := &server{http.NewServeMux(), client}
	// api
	s.HandleFunc("/create", s.create)
	s.HandleFunc("/view", s.view)
	s.HandleFunc("/subs", s.subs)
	s.HandleFunc("/add", s.add)
	s.HandleFunc("/remove", s.remove)
	s.HandleFunc("/send", s.send)

	// views
	s.HandleFunc("/audience", s.viewAudience)
	s.HandleFunc("/", s.viewAudiences)
	return s
}

// server composes http's ServeMux struct for http routing
type server struct {
	*http.ServeMux
	client *client
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	aud, _ := s.client.getAudience(r.URL.Query().Get("aud"))
	ctx := context.WithValue(r.Context(), "aud", aud)
	s.ServeMux.ServeHTTP(w, r.WithContext(ctx))
}

// create creates a new aud
// [POST] /create
func (s server) create(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.restrictMethods("POST")
	h.respond(db.saveAudience(h.aud()))
}

// view gets basic details for a aud
// [GET] /view
func (s server) view(w http.ResponseWriter, r *http.Request) {
	h, _ := s.handle(w, r)
	defer h.catch()

	h.restrictMethods(http.MethodGet)
	// h.respond(h.aud())
}

// subs gets all subscriptsion for a aud
// [GET] /subs
func (s server) subs(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.restrictMethods("GET")
	h.respond(db.getMembers(h.aud()))
}

// add creates a Member for a aud
// [POST] /add { listName, name, address }
func (s server) add(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.restrictMethods("POST")
	h.respond(db.saveMember(&Member{
		aud:   h.aud(),
		Name:  r.FormValue("name"),
		Email: r.FormValue("address"),
	}))
}

// remove removes a Member from a aud by address
// [DELETE] /add { listName, address }
func (s server) remove(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.restrictMethods("DELETE")

	h.respond(nil, db.deleteMember(
		h.aud(),
		r.FormValue("address"),
	))
}

// send sends mail to all Members on a aud
// [POST] /send { message, listName }
func (s server) send(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.restrictMethods("POST")
	subs, err := db.getMembers(h.aud())
	h.respond(nil, err, h.gmail().SendEmail(h.body(), subs...))
}

// handle creates a new handler, we also defer a recovery func
// to handle serverErrors we encounter
func (s server) handle(w http.ResponseWriter, r *http.Request) (*handler, *client) {
	return &handler{w, r, nil, nil, 200}, s.client
}
