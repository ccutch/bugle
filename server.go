package bugle

import (
	"context"
	"net/http"
)

// Server creates a new router and binds all routes
func Server(db *database) http.Handler {
	s := &server{http.NewServeMux(), db}
	s.staticFiles("public")
	s.HandleFunc("/audience", s.viewAudience)
	s.HandleFunc("/", s.viewAudiences)
	return s
}

// server composes http's ServeMux struct for http routing
type server struct {
	*http.ServeMux
	db *database
}

// ServeHTTP for http.Handler interface. We are loading variables
// into context and use underlying serve mux with new context into
func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	aud, _ := s.db.getAudience(ctx, r.URL.Query().Get("aud"))

	ctx = context.WithValue(ctx, "aud", aud)
	ctx = context.WithValue(ctx, "user", User("connor@bugl.email", "...token..."))
	ctx = context.WithValue(ctx, "api", r.URL.Query().Get("api") == "true")

	s.ServeMux.ServeHTTP(w, r.WithContext(ctx))
}

// handle creates a new handler, we also defer a recovery func
// to handle serverErrors we encounter
func (s server) handle(w http.ResponseWriter, r *http.Request) (*handler, *database) {
	ctx, cancel := context.WithCancel(r.Context())
	return &handler{w, r.WithContext(ctx), nil, ctx, cancel, nil, 200}, s.db
}

// staticFiles sets up file server and handles http requests
// with the same name
func (s server) staticFiles(d string) {
	fs := http.FileServer(http.Dir("./" + d))
	s.Handle("/"+d+"/", http.StripPrefix("/"+d+"/", fs))
}
