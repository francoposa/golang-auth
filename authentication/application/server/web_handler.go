package server

import (
	"fmt"
	"net/http"
)

type AuthNWebHandler struct {
	templateRenderer     *TemplateRenderer
	loginTemplateName    string
	registerTemplateName string
}

func NewWebHandler(
	templateRenderer *TemplateRenderer,
	loginTemplateName string,
	registerTemplateName string,
) *AuthNWebHandler {
	return &AuthNWebHandler{
		templateRenderer:     templateRenderer,
		loginTemplateName:    loginTemplateName,
		registerTemplateName: registerTemplateName,
	}
}

func (h *AuthNWebHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	err := h.templateRenderer.RenderTemplate(w, h.loginTemplateName, "")
	fmt.Println(err)
}

func (h *AuthNWebHandler) GetRegister(w http.ResponseWriter, r *http.Request) {
	err := h.templateRenderer.RenderTemplate(w, h.registerTemplateName, "")
	fmt.Println(err)
}
