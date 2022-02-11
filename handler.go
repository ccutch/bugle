package bugle

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
)

var globalTemplates []string = []string{
	"templates/layouts/main.html",
}

// request is for an individual request
type handler struct {
	w http.ResponseWriter
	r *http.Request
	t *template.Template

	err  error
	code int
}

// listName is getter for url querystring value "aud"
func (h handler) aud() *Audience {
	if aud, ok := h.r.Context().Value("aud").(Audience); ok {
		return &aud
	}
	return nil
}

// gmail is a getter for gmail client
func (h handler) gmail() *GmailClient {
	return Gmail(h.r.Context())
}

// body is a simple body parser
func (h handler) body() string {
	body, err := ioutil.ReadAll(h.r.Body)
	h.handle(err, http.StatusBadRequest)
	return string(body)
}

// fail error
func (h handler) handle(err error, code int) {
	if err != nil {
		panic(serverError{err, code})
	}
}

// serverError to encapsolate state
type serverError struct {
	err  error
	code int
}

func (h handler) restrictMethods(methods ...string) {
	for _, m := range methods {
		if h.r.Method == m {
			return
		}
	}
	h.handle(errors.New("Invalid method"), http.StatusMethodNotAllowed)
}

func (h handler) requireAudience(methods ...string) {
	if aud := h.aud(); aud == nil {
		h.handle(errors.New("Audience query param required"), http.StatusMethodNotAllowed)
	}
}

// respond is a helper function for normalize errorhandle
func (h handler) respond(v interface{}, errs ...error) {
	h.w.WriteHeader(h.code)
	switch errs = clean(errs); { // I wrote this way because worse with ifs
	case len(errs) > 0:
		buff := bytes.NewBuffer([]byte("Error:\n\n"))
		for _, err := range errs {
			buff.WriteString(err.Error() + "\n")
		}
		buff.WriteTo(h.w)

	case h.t != nil:
		h.w.Header().Add("Content-Type", "text/html")
		h.handle(h.t.Execute(h.w, v), http.StatusInternalServerError)

	default:
		h.handle(
			json.NewEncoder(h.w).Encode(v),
			http.StatusInternalServerError,
		)
	}
}

func clean(errs []error) (res []error) {
	for _, err := range errs {
		if err != nil {
			res = append(res, err)
		}
	}
	return res
}

func (h *handler) loadView(name string) {
	h.t, h.err = template.ParseFiles("views/" + name + ".html")
	if h.err != nil {
		h.code = http.StatusNotFound
	}
}

func (h handler) catch() {
	if r, ok := recover().(serverError); ok {
		h.code = r.code
		h.respond(nil, r.err)
	}
}
