package bugle

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	*http.ServeMux
	db    DBClient
	gmail GmailClient
}

func NewServer(db DBClient, gmail GmailClient) *Server {
	s := &Server{http.NewServeMux(), db, gmail}
	s.HandleFunc("/create", s.CreateList)
	s.HandleFunc("/view", s.ViewList)
	s.HandleFunc("/add", s.AddEmail)
	s.HandleFunc("/remove", s.RemoveEmail)
	s.HandleFunc("/send", s.SendMailBlast)
	return s
}

func (s *Server) CreateList(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ListName string `json:"listName"`
	}

	json.NewDecoder(r.Body).Decode(&body)

}

func (s *Server) ViewList(w http.ResponseWriter, r *http.Request) {}

func (s *Server) AddEmail(w http.ResponseWriter, r *http.Request) {}

func (s *Server) RemoveEmail(w http.ResponseWriter, r *http.Request) {}

func (s *Server) SendMailBlast(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Message  string `json:"message"`
		ListName string `json:"listName"`
	}

	json.NewDecoder(r.Body).Decode(&body)
	emails, _ := s.db.LookupEmails(body.ListName)
	s.gmail.SendEmail(body.Message, emails...)
}
