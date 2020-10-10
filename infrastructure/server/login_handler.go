package server

import (
	"fmt"
	"net/http"

	"golang-auth/usecases/repos"
)

type LoginHandler struct {
	authNUserRepo     repos.AuthNUserRepo
	templateRenderer  *TemplateRenderer
	loginTemplateName string
}

func NewLoginHandler(
	authNUserRepo repos.AuthNUserRepo,
	templateRenderer *TemplateRenderer,
	loginTemplateName string,
) *LoginHandler {
	return &LoginHandler{
		authNUserRepo:     authNUserRepo,
		templateRenderer:  templateRenderer,
		loginTemplateName: loginTemplateName,
	}
}

func (h *LoginHandler) Get(w http.ResponseWriter, r *http.Request) {
	err := h.templateRenderer.RenderTemplate(w, h.loginTemplateName, "")
	fmt.Println(err)
}
