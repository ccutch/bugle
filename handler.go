package bugle

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// request is for an individual request
type handler struct {
	w http.ResponseWriter
	r *http.Request
}

// listName is getter for url querystring value "list"
func (h handler) listName() string {
	return h.r.Context().Value("list").(string)
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

// respond is a helper function for normalize errorhandle
func (h handler) respond(v interface{}, errs ...error) {
	switch {
	case len(errs) > 0:
		h.w.WriteHeader(http.StatusInternalServerError)

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
