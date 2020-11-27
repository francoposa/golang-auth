package server

import (
	"fmt"
	"net/http"
)

type WebHandler struct {
	templateRenderer     *TemplateRenderer
	loginTemplateName    string
	registerTemplateName string
}

func NewWebHandler(
	templateRenderer *TemplateRenderer,
	loginTemplateName string,
	registerTemplateName string,
) *WebHandler {
	return &WebHandler{
		templateRenderer:     templateRenderer,
		loginTemplateName:    loginTemplateName,
		registerTemplateName: registerTemplateName,
	}
}

func (h *WebHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	err := h.templateRenderer.RenderTemplate(w, h.loginTemplateName, "")
	fmt.Println(err)
}

func (h *WebHandler) GetRegister(w http.ResponseWriter, r *http.Request) {
	err := h.templateRenderer.RenderTemplate(w, h.registerTemplateName, "")
	fmt.Println(err)
}
