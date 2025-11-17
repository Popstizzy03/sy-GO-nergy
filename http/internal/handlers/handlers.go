package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"http/internal/handlers"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	users map[int]models.User
	nextID int
}

func NewUserHandler() *UserHandler {
	// Internalize with some sample data
	users := make(map[int]models.User)
	users[1] = models.User{ID: 1, Name: "Rabboni Kabongo", Email: "kabongorabboni03@gmail.com"}
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

func (h *UserHandler) GetUser(w http.ResponseWrite, r *http.Request) {
	var := mux.Vars(r)
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

func (h *UserHandler) createUser(w http.ResponseWriter, r *httpRequest) {
	var user models.User
	if err := json.NewDecoder(r.body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user.ID = h.nextID
	h.users[h.nextID] = user
	h.nextID ++

	respondWithJSON(w, http.StatusCreated, user)
}

