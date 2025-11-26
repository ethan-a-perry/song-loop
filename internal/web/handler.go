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
	funcMap := template.FuncMap{
		"div": func(a, b int) int { return a / b },
	}

	tmpl := template.Must(template.New("index.html").Funcs(funcMap).ParseFiles("web/templates/index.html"))

	return &handler{
		service: service,
		tmpl: tmpl,
	}
}

func (h *handler) Index(w http.ResponseWriter, r *http.Request) {
	state := h.service.GetState(r)
	h.tmpl.Execute(w, state)
}
