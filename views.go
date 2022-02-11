package bugle

import "net/http"

func (s server) viewAudiences(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.restrictMethods("GET", "POST")
	h.loadView("audience-list")

	h.respond(db.getAudienceForUser("Connor"))
}

func (s server) viewAudience(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.restrictMethods("GET", "POST")
	h.requireAudience()
	h.loadView("audience")

	m, err := db.getMembers(h.aud())
	h.respond(&struct {
		Audience *Audience `json:"audience"`
		Members  []Member  `json:"members"`
	}{h.aud(), m}, err)
}
