package bugle

import "net/http"

type user struct {
	Token string `json:"-"`
	Name  string
	Email string
}

func User(email, token string) (u user, err error) {
	u.Email = email
	u.Token = token

	// load name from google api
	return u, err
}

func parseUser(r *http.Request) (u user, err error) {
	u.Token = r.Header.Get("x-goog-iap-jwt-assertion")
	u.Email = r.Header.Get("x-goog-authenticated-user-email")

	// TODO remove after testing
	if u.Email == "" {
		u.Email = "connor@bugl.email"
	}

	return u, err
}
