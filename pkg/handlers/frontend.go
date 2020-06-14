package handlers

import (
	"html/template"
	"net/http"
)

type FrontendHandler struct{}

func NewFrontendHandler() *FrontendHandler {
	return &FrontendHandler{}
}

func (h *FrontendHandler) IndexPage(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("./template/index.html"))
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, `{"message":"template error"}`, http.StatusInternalServerError)
	}
}
