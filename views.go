package bugle

import (
	"net/http"
)

func (s server) viewAudiences(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.requireUser()
	h.restrictMethods("GET", "POST")
	h.loadView("dashboard/audience-list", "dashboard/_layout")

	auds, err := db.getAudienceForUser(h.ctx, h.user())
	h.respond(&struct {
		User      *user
		Audiences []Audience
	}{h.user(), auds}, err)
}

// viewAudience [GET | POST] /audience
// Methods for audience, viewing and adding member
func (s server) viewAudience(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.requireUser()
	h.requireAudience()
	h.restrictMethods("GET", "POST")
	h.loadView("dashboard/audience", "dashboard/_layout")

	mems, err := db.getMembers(h.ctx, h.aud())
	h.respond(&struct {
		Audience *Audience
		Members  []Member
	}{h.aud(), mems}, err)
}
