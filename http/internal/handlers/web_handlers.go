package handlers

import (
	"html/template"
	"http/internal/models"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type WebHandler struct {
	userHandler *UserHandler
	templates   *template.Template
}

func NewWebHandler(userHandler *UserHandler) *WebHandler {
	wh := &WebHandler{
		userHandler: userHandler,
	}
	
	// Load templates immediately
	err := wh.loadTemplates()
	if err != nil {
		log.Printf("Warning: Could not load templates: %v", err)
	}
	
	return wh
}

func (wh *WebHandler) loadTemplates() error {
	// Define the path to the templates directory
	templatesPath := "web/templates/*.html"

	// Log the path for debugging
	log.Printf("Loading templates from: %s", templatesPath)

	// Parse all template files in the directory
	tmpl, err := template.ParseGlob(templatesPath)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		return err
	}

	// Assign the parsed templates to the handler
	wh.templates = tmpl
	log.Println("Successfully loaded templates")
	return nil
}

func (wh *WebHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	log.Printf("HomePage handler called")
	wh.renderTemplate(w, "index.html", nil)
}

func (wh *WebHandler) UsersPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("UsersPage handler called")
	
	// Get users from the user handler
	users := make([]models.User, 0, len(wh.userHandler.users))
	for _, user := range wh.userHandler.users {
		users = append(users, user)
	}
	
	data := struct {
		Users []models.User
	}{
		Users: users,
	}
	
	wh.renderTemplate(w, "users.html", data)
}

func (wh *WebHandler) CreateUserPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateUserPage handler called")
	wh.renderTemplate(w, "create-user.html", nil)
}

func (wh *WebHandler) renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	log.Printf("Attempting to render template: %s", tmpl)
	
	if wh.templates == nil {
		log.Printf("ERROR: Templates not loaded!")
		http.Error(w, "Templates not available", http.StatusInternalServerError)
		return
	}
	
	// Set HTML content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	// Execute the "base" template, which will in turn render the content of the specific template
	err := wh.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Error rendering template %s: %v", tmpl, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	log.Printf("Successfully rendered template: %s", tmpl)
}

func (wh *WebHandler) ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	// Security: Prevent directory traversal
	if r.URL.Path == "/static/" {
		http.NotFound(w, r)
		return
	}
	
	// Strip "/static" prefix and serve from web/static directory
	path := r.URL.Path[len("/static/"):]
	
	// Try different possible static file locations
	possiblePaths := []string{
		filepath.Join("web/static", path),
		filepath.Join("./web/static", path),
		filepath.Join("../web/static", path),
	}
	
	for _, filePath := range possiblePaths {
		if _, err := os.Stat(filePath); err == nil {
			log.Printf("Serving static file: %s", filePath)
			
			// Set caching headers for static assets
			w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour
			http.ServeFile(w, r, filePath)
			return
		}
	}
	
	// File not found
	log.Printf("Static file not found: %s", path)
	http.NotFound(w, r)
}
