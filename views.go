package bugle

import (
	"net/http"
	"time"
)

// viewAudience is a controller in bugle's server
//
// [GET]  => fetch all audiences for current user
// [POST] => create new audience for current user
func (s server) viewAudiences(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.restrictMethods("GET", "POST")

	switch r.Method {
	case "GET":
		auds, err := db.getAudienceForUser(h.ctx, h.user())
		h.handle(err, http.StatusInternalServerError)
		h.loadView("dashboard/audience-list", "dashboard/_layout")
		h.respond(&struct {
			User      *user
			Audiences []Audience
		}{h.user(), auds})

	case "POST":
		aud := Audience{nil, r.FormValue("name"), h.user().Email, time.Now()}
		mem := Member{nil, &aud, "New member", h.user().Email, time.Now()}
		h.handle(db.saveAudience(h.ctx, &aud), http.StatusConflict)
		h.handle(db.saveMember(h.ctx, &mem), http.StatusInternalServerError)
		h.redirect("/audience?aud=" + aud.KeyName())
	}
}

// viewAudience [GET | POST] /audience
// Methods for audience, viewing and adding member
func (s server) viewAudience(w http.ResponseWriter, r *http.Request) {
	h, db := s.handle(w, r)
	h.requireAudience()
	h.restrictMethods("GET", "POST")

	h.loadView("dashboard/audience", "dashboard/_layout")
	mems, err := db.getMembers(h.ctx, h.aud())
	h.respond(&struct {
		Audience *Audience
		Members  []Member
	}{h.aud(), mems}, err)
}
