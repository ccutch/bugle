package bugle

import (
	"context"
	"net/http"
)

// Server creates a new router and binds all routes
func Server(db *database) http.Handler {
	s := &server{http.NewServeMux(), db}
	s.static("public")
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
	// using and updating request context per request
	var ctx = r.Context()
	// load audience & user, ignoring error for zero value
	var aud, _ = s.db.getAudience(ctx, r.URL.Query().Get("aud"))
	var usr, _ = parseUser(r)
	// context variables
	ctx = context.WithValue(ctx, "usr", usr)
	ctx = context.WithValue(ctx, "aud", aud)
	// context flags
	ctx = context.WithValue(ctx, "api", r.URL.Query().Get("api") == "true")
	// defer serve http to serve mux
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
func (s server) static(d string) {
	fs := http.FileServer(http.Dir("./" + d))
	s.Handle("/"+d+"/", http.StripPrefix("/"+d+"/", fs))
}
