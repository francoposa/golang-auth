package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func NewTemplates(pattern, baseTemplatePath string) map[string]*template.Template {
	templates := make(map[string]*template.Template)

	templatePaths, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}

	for _, templatePath := range templatePaths {
		if templatePath == baseTemplatePath {
			continue
		}
		fileName := filepath.Base(templatePath)
		templates[fileName] = template.Must(
			template.ParseFiles(templatePath, baseTemplatePath),
		)
	}

	return templates
}

type TemplateRenderer struct {
	templates        map[string]*template.Template
	baseTemplatePath string
}

func NewTemplateRenderer(templates map[string]*template.Template, baseTemplateName string) *TemplateRenderer {
	return &TemplateRenderer{templates: templates, baseTemplatePath: baseTemplateName}
}

func (tr *TemplateRenderer) RenderTemplate(w http.ResponseWriter, templateFileName string, data interface{}) error {
	templateToRender, ok := tr.templates[templateFileName]
	if !ok {
		return fmt.Errorf("The templateToRender %s does not exist.", templateFileName)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return templateToRender.ExecuteTemplate(w, tr.baseTemplatePath, data)
}
