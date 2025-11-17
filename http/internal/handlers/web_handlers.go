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
	// Try different possible template locations
	possiblePaths := []string{
		"web/templates/*.html",
		"./web/templates/*.html",
		"../web/templates/*.html",
	}
	
	var tmpl *template.Template
	var err error
	
	for _, path := range possiblePaths {
		log.Printf("Trying to load templates from: %s", path)
		
		// Check if any files match this pattern
		matches, _ := filepath.Glob(path)
		if len(matches) > 0 {
			log.Printf("Found template files: %v", matches)
			tmpl, err = template.ParseGlob(path)
			if err == nil {
				wh.templates = tmpl
				log.Printf("Successfully loaded templates from: %s", path)
				return nil
			}
		}
	}
	
	// If we get here, try manual loading as fallback
	log.Printf("Trying manual template loading...")
	return wh.loadTemplatesManually()
}

func (wh *WebHandler) loadTemplatesManually() error {
	// List all template files we expect
	templateFiles := []string{
		"base.html",
		"index.html", 
		"users.html",
		"create-user.html",
	}
	
	tmpl := template.New("").Funcs(template.FuncMap{
		"upper": func(s string) string {
			if len(s) > 0 {
				return string(s[0])
			}
			return s
		},
	})
	
	basePath := "web/templates/"
	
	// Try to parse each template file individually
	for _, filename := range templateFiles {
		filePath := filepath.Join(basePath, filename)
		
		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Printf("Template file does not exist: %s", filePath)
			continue
		}
		
		// Read and parse the template
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading template %s: %v", filename, err)
			continue
		}
		
		_, err = tmpl.Parse(string(content))
		if err != nil {
			log.Printf("Error parsing template %s: %v", filename, err)
			return err
		}
		
		log.Printf("Successfully loaded template: %s", filename)
	}
	
	wh.templates = tmpl
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
	
	// Execute the template
	err := wh.templates.ExecuteTemplate(w, tmpl, data)
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
