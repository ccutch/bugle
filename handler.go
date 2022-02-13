package bugle

import (
	"bytes"
	"context"
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

	ctx    context.Context
	cancel context.CancelFunc

	err  error
	code int
}

// listName is getter for url querystring value "aud"
func (h handler) aud() *Audience {
	if aud, ok := h.ctx.Value("aud").(Audience); ok {
		return &aud
	}
	return nil
}

func (h handler) api() bool {
	api, _ := h.ctx.Value("api").(bool)
	return api
}

func (h handler) user() *user {
	if user, ok := h.ctx.Value("user").(user); ok {
		return &user
	}
	return nil
}

// gmail is a getter for gmail client
func (h handler) gmail() *GmailClient {
	return Gmail(h.ctx)
}

// body is a simple body parser
func (h *handler) body() string {
	body, err := ioutil.ReadAll(h.r.Body)
	h.handle(err, http.StatusBadRequest)
	return string(body)
}

// fail error
func (h *handler) handle(err error, code int) {
	if err == nil {
		return
	}

	h.err = err
	h.code = code
	h.cancel()
}

func (h *handler) restrictMethods(methods ...string) {
	for _, m := range methods {
		if h.r.Method == m {
			return
		}
	}

	h.handle(errors.New("Invalid method"), http.StatusMethodNotAllowed)
}

func (h *handler) requireAudience() {
	aud := h.aud()
	if aud.IsZero() {
		err := errors.New("Audience query param required")
		h.handle(err, http.StatusBadRequest)
	}
}

func (h *handler) requireUser() {
	if user := h.user(); user == nil {
		err := errors.New("User requried")
		h.handle(err, http.StatusBadRequest)
	}
}

// respond is a helper function for normalize errorhandle
func (h *handler) respond(v interface{}, errs ...error) {
	h.w.WriteHeader(h.code)
	switch errs = clean(append(errs, h.err)); { // I wrote this way because worse with ifs

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

func (h *handler) loadView(names ...string) {
	if h.api() {
		return
	}

	var files []string
	for _, name := range names {
		files = append(files, "views/"+name+".html")
	}

	var err error
	h.t, err = template.ParseFiles(files...)
	if err != nil {
		h.handle(err, http.StatusNotFound)
	}
}
