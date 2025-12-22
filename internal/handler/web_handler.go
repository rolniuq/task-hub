package handler

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type WebHandler struct {
	templates *template.Template
}

func NewWebHandler(templatesDir string) (*WebHandler, error) {
	tmpl, err := template.ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return nil, err
	}

	return &WebHandler{
		templates: tmpl,
	}, nil
}

func (h *WebHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.render(w, "login.html", nil)
}

func (h *WebHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.render(w, "register.html", nil)
}

func (h *WebHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.render(w, "dashboard.html", nil)
}

func (h *WebHandler) render(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
