package bugle

import "net/http"

type user struct {
	Name  string
	Email string
	token string
}

func parseUser(r *http.Request) (u user, err error) {
	if e := r.Header.Get("x-goog-authenticated-user-email"); e == "" {
		u.Email = "test@bugl.email"
	} else {
		u.Email = e[len("accounts.google.com:"):]
	}
	u.token = r.Header.Get("x-goog-iap-jwt-assertion")
	return u, err
}
