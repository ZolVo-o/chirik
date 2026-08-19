package handlers

import (
	"chirik/internal/middleware"
	"chirik/internal/repository"
	"net/http"
)

type UsersHandler struct {
	repo *repository.Repository
}

func NewUsersHandler(repo *repository.Repository) *UsersHandler {
	return &UsersHandler{repo: repo}
}

func (h *UsersHandler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := middleware.GetUserID(r)
	if !ok {
		sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	users, err := h.repo.SearchUsers(r.URL.Query().Get("q"), userID)
	if err != nil {
		sendJSONError(w, "Internal error", http.StatusInternalServerError)
		return
	}
	sendJSON(w, users)
}
