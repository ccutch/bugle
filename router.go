package bugle

import (
	"net/http"
)

// Router composes http's ServeMux struct for http routing
type Router struct{ *http.ServeMux }

// NewRouter creates a new router and binds all routes
func NewRouter() *Router {
	s := &Router{http.NewServeMux()}
	s.HandleFunc("/create", s.CreateList)
	s.HandleFunc("/view", s.ViewList)
	s.HandleFunc("/add", s.AddEmail)
	s.HandleFunc("/remove", s.RemoveEmail)
	s.HandleFunc("/send", s.SendMailBlast)
	return s
}

// CreateList creates a new list
// [POST] /create { listName }
func (Router) CreateList(w http.ResponseWriter, r *http.Request) {
	body, db := parse(r.Body), NewDBClient(r.Context())

	emails := db.NewList(body.ListName)
	respond(w, emails, db.err)
}

// ViewList gets all subscriptsion for a list
// [GET] /view?list=my%20friends
func (Router) ViewList(w http.ResponseWriter, r *http.Request) {
	listName, db := r.URL.Query().Get("list"), NewDBClient(r.Context())

	emails := db.GetSubscriptions(listName)
	respond(w, emails, db.err)
}

// AddEmail creates a subscription for a list
// [POST] /add { listName, name, address }
func (Router) AddEmail(w http.ResponseWriter, r *http.Request) {
	body, db := parse(r.Body), NewDBClient(r.Context())

	subscription := db.NewSubscription(body.ListName, body.Name, body.Address)
	respond(w, subscription, db.err)
}

// RemoveEmail removes a subscription from a list by address
// [POST] /add { listName, address }
func (Router) RemoveEmail(w http.ResponseWriter, r *http.Request) {
	body, db := parse(r.Body), NewDBClient(r.Context())

	db.DeleteSubscription(body.ListName, body.Address)
	respond(w, nil, db.err)
}

// SendMailBlast sends mail to all subscriptions on a list
// [POST] /send { message, listName }
func (Router) SendMailBlast(w http.ResponseWriter, r *http.Request) {
	body, db, gmail := parse(r.Body), NewDBClient(r.Context()), NewGmailClient(r.Context())

	emails := db.GetSubscriptions(body.ListName)
	err := gmail.SendEmail(body.Message, emails...)
	respond(w, nil, db.err, err)
}
