package bugle

import (
	"encoding/json"
	"net/http"
)

type Router struct{ *http.ServeMux }

func NewRouter() *Router {
	s := &Router{http.NewServeMux()}
	s.HandleFunc("/create", s.CreateList)
	s.HandleFunc("/view", s.ViewList)
	s.HandleFunc("/add", s.AddEmail)
	s.HandleFunc("/remove", s.RemoveEmail)
	s.HandleFunc("/send", s.SendMailBlast)
	return s
}

func (s *Router) CreateList(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ListName string `json:"listName"`
	}

	json.NewDecoder(r.Body).Decode(&body)
	db := NewDBClient(r.Context())
	emails := db.NewList(body.ListName)
	json.NewEncoder(w).Encode(emails)
}

func (s *Router) ViewList(w http.ResponseWriter, r *http.Request) {
	listName := r.URL.Query().Get("list")

	db := NewDBClient(r.Context())
	emails := db.GetSubscriptions(listName)
	json.NewEncoder(w).Encode(emails)
}

func (s *Router) AddEmail(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ListName string `json:"listName"`
		Name     string `json:"name"`
		Address  string `json:"address"`
	}

	json.NewDecoder(r.Body).Decode(&body)
	db := NewDBClient(r.Context())
	subscription := db.NewSubscription(body.ListName, body.Name, body.Address)
	json.NewEncoder(w).Encode(subscription)
}

func (s *Router) RemoveEmail(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ListName string `json:"listName"`
		Address  string `json:"address"`
	}

	json.NewDecoder(r.Body).Decode(&body)
	db := NewDBClient(r.Context())
	db.DeleteSubscription(body.ListName, body.Address)
	w.WriteHeader(200)
}

func (s *Router) SendMailBlast(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Message  string `json:"message"`
		ListName string `json:"listName"`
	}

	json.NewDecoder(r.Body).Decode(&body)
	db := NewDBClient(r.Context())
	gmail := NewGmailClient(r.Context())
	emails := db.GetSubscriptions(body.ListName)
	gmail.SendEmail(body.Message, emails...)
	w.WriteHeader(200)
}
