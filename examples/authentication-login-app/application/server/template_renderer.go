package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
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
		fileNameNoExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		templates[fileNameNoExt] = template.Must(
			template.ParseFiles(templatePath, baseTemplatePath),
		)
	}

	return templates
}

type TemplateRenderer struct {
	templates map[string]*template.Template
}

func NewTemplateRenderer(templates map[string]*template.Template) *TemplateRenderer {
	return &TemplateRenderer{templates: templates}
}

func (tr *TemplateRenderer) RenderTemplate(w http.ResponseWriter, templateName string, data interface{}) error {
	templateToRender, ok := tr.templates[templateName]
	if !ok {
		return fmt.Errorf("Template %s does not exist.", templateName)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return templateToRender.ExecuteTemplate(w, templateName, data)
}
