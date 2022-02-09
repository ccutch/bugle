package bugle

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// request is for an individual request
type handler struct {
	w    http.ResponseWriter
	r    *http.Request
	err  error
	code int
}

// listName is getter for url querystring value "list"
func (h handler) list() *MailingList {
	list := h.r.Context().Value("list").(MailingList)
	return &list
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
func (h *handler) handle(err error, code int) {
	if err != nil {
		h.code = code
		h.err = err
	}
}

// serverError to encapsolate state
type serverError struct {
	err  error
	code int
}

func (h *handler) restrictMethods(methods ...string) {
	for _, m := range methods {
		if h.r.Method == m {
			return
		}
	}

	h.err = errors.New("Invalid method")
	h.code = http.StatusBadRequest
}

// respond is a helper function for normalize errorhandle
func (h handler) respond(v interface{}, errs ...error) {
	errs = clean(errs...)

	switch {
	case len(errs) > 0:
		if h.code > 0 {
			h.w.WriteHeader(h.code)
		} else {
			h.w.WriteHeader(http.StatusInternalServerError)
		}

		var buff bytes.Buffer
		buff.WriteString("Error:\n\n")

		for _, err := range errs {
			buff.WriteString(err.Error() + "\n")
		}
		buff.WriteTo(h.w)

	default:
		err := json.NewEncoder(h.w).Encode(v)
		h.handle(err, http.StatusInternalServerError)
	}
}

func clean(errs ...error) (res []error) {
	for _, err := range errs {
		if err != nil {
			res = append(res, err)
		}
	}
	return res
}
