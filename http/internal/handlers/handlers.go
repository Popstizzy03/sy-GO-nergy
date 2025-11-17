package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"http/internal/models"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	users map[int]models.User
	nextID int
}

func NewUserHandler() *UserHandler {
	// Initialize with some sample data
	users := make(map[int]models.User)
	users[1] = models.User{ID: 1, Name: "John Doe", Email: "john@example.com"}
	users[2] = models.User{ID: 2, Name: "Jane Smith", Email: "jane@example.com"}
	
	return &UserHandler{
		users: users,
		nextID: 3,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := make([]models.User, 0, len(h.users))
	for _, user := range h.users {
		users = append(users, user)
	}
	
	respondWithJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	
	user, exists := h.users[id]
	if !exists {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	
	respondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	
	user.ID = h.nextID
	h.users[h.nextID] = user
	h.nextID++
	
	respondWithJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	}
	respondWithJSON(w, http.StatusOK, response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	userHandler := NewUserHandler()
	webHandler := NewWebHandler(userHandler)

	// Web routes
	router.HandleFunc("/", webHandler.HomePage).Methods("GET")
	router.HandleFunc("/users", webHandler.UsersPage).Methods("GET")
	router.HandleFunc("/users/create", webHandler.CreateUserPage).Methods("GET")
	
	// Serve static files
	staticFileDirectory := http.Dir("./web/static/")
	staticFileServer := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))
	router.PathPrefix("/static/").Handler(staticFileServer)
	
	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUser).Methods("GET")
	
	// Health check
	router.HandleFunc("/health", userHandler.HealthCheck).Methods("GET")
	
	// Root endpoint
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Go HTTP Server is running!",
			"version": "1.0.0",
		})
	}).Methods("GET")
	
	return router
}
