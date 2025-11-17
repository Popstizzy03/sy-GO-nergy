package handlers

import (
	"html/template"
	"http/internal/models"
	"log"
	"net/http"
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


