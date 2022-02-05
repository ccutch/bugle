package bugle

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type validInputs struct {
	ListName string `json:"listName"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Message  string `json:"message"`
}

func parse(r io.Reader) (inputs validInputs) {
	json.NewDecoder(r).Decode(&inputs)
	return inputs
}

func respond(w http.ResponseWriter, v interface{}, errs ...error) {
	switch {

	case len(errs) > 0:
		w.WriteHeader(500)

		var buff bytes.Buffer
		buff.WriteString("Error:\n\n")

		for _, err := range errs {
			buff.WriteString(err.Error() + "\n")
		}
		buff.WriteTo(w)

	default:
		json.NewEncoder(w).Encode(v)
	}
}
