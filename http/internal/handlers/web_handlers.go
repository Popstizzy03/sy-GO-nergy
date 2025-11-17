package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"http/internal/models"
)

type WebHandler struct {
	userHandler *UserHandler
	templates   *template.Template
}

func NewWebHandler(userHandler *UserHandler) *WebHandler {
	wh := &WebHandler{
		userHandler: userHandler,
	}
	wh.loadTemplates()
	return wh
}

func (wh *WebHandler) loadTemplates() {
	tmpl := template.New("").Funcs(template.FuncMap{
		"upper": func(s string) string {
			if len(s) > 0 {
				return string(s[0])
			}
			return s
		},
	})
	
	templatesPath := "web/templates/*.html"
	tmpl, err := tmpl.ParseGlob(templatesPath)
	if err != nil {
		panic("Failed to load templates: " + err.Error())
	}
	
	wh.templates = tmpl
}

func (wh *WebHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	wh.renderTemplate(w, "index.html", nil)
}

func (wh *WebHandler) UsersPage(w http.ResponseWriter, r *http.Request) {
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
	wh.renderTemplate(w, "create-user.html", nil)
}

func (wh *WebHandler) renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := wh.templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ServeStaticFiles handles serving CSS, JS, and images
func (wh *WebHandler) ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	// Security: Prevent directory traversal
	if r.URL.Path == "/static/" {
		http.NotFound(w, r)
		return
	}
	
	// Strip "/static" prefix and serve from web/static directory
	path := r.URL.Path[len("/static/"):]
	filePath := filepath.Join("web/static", path)
	
	// Set caching headers for static assets
	w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour
	http.ServeFile(w, r, filePath)
}
