package web

import (
	"html/template"
	"net/http"
)

type handler struct {
	service *Service
	tmpl *template.Template
}

func NewHandler(service *Service) *handler {
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))

	return &handler{
		service: service,
		tmpl: tmpl,
	}
}

func (h *handler) Index(w http.ResponseWriter, r *http.Request) {
	state := h.service.GetState(r)
	h.tmpl.Execute(w, state)
}
