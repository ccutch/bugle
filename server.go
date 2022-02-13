package bugle

import (
	"context"
	"net/http"
)

// Server creates a new router and binds all routes
func Server(client *client) http.Handler {
	s := &server{http.NewServeMux(), client}
	s.staticFiles("public")
	s.HandleFunc("/audience", s.viewAudience)
	s.HandleFunc("/", s.viewAudiences)
	return s
}

// server composes http's ServeMux struct for http routing
type server struct {
	*http.ServeMux
	client *client
}

// ServeHTTP for http.Handler interface. We are loading variables
// into context and use underlying serve mux with new context into
func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	aud, _ := s.client.getAudience(r.URL.Query().Get("aud"))
	ctx := context.WithValue(r.Context(), "aud", aud)
	ctx = context.WithValue(ctx, "user", User("connor@bugl.email", "...token..."))
	s.ServeMux.ServeHTTP(w, r.WithContext(ctx))
}

// handle creates a new handler, we also defer a recovery func
// to handle serverErrors we encounter
func (s server) handle(w http.ResponseWriter, r *http.Request) (*handler, *client) {
	api := r.URL.Query().Get("api")
	return &handler{w, r, nil, api == "true", nil, 200}, s.client
}

// staticFiles sets up file server and handles http requests
// with the same name
func (s server) staticFiles(d string) {
	fs := http.FileServer(http.Dir("./" + d))
	s.Handle("/"+d+"/", http.StripPrefix("/"+d+"/", fs))
}
