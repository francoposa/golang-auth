package server

import (
	"fmt"
	"net/http"
)

type WebHandler struct {
	templateRenderer *TemplateRenderer
}

func NewWebHandler(
	templateRenderer *TemplateRenderer,
) *WebHandler {
	return &WebHandler{
		templateRenderer: templateRenderer,
	}
}

func (h *WebHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	err := h.templateRenderer.RenderTemplate(w, "sign-in", "")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (h *WebHandler) GetRegister(w http.ResponseWriter, r *http.Request) {
	err := h.templateRenderer.RenderTemplate(w, "sign-up", "")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
