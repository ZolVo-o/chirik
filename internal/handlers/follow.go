package handlers

import (
	"chirik/internal/middleware"
	"chirik/internal/realtime"
	"chirik/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type FollowHandler struct {
	repo      *repository.Repository
	publisher realtime.Publisher
}

func NewFollowHandler(repo *repository.Repository, publisher realtime.Publisher) *FollowHandler {
	return &FollowHandler{repo: repo, publisher: publisher}
}

func (h *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	followerID, ok := middleware.GetUserID(r)
	if !ok {
		h.sendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		h.sendError(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	followingID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		h.sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if followerID == followingID {
		h.sendError(w, "Cannot follow yourself", http.StatusBadRequest)
		return
	}
	user, err := h.repo.GetUserByID(followingID)
	if err != nil {
		h.sendError(w, "Internal error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		h.sendError(w, "User not found", http.StatusNotFound)
		return
	}
	alreadyFollowing, err := h.repo.IsFollowing(followerID, followingID)
	if err != nil {
		h.sendError(w, "Internal error", http.StatusInternalServerError)
		return
	}
	if alreadyFollowing {
		h.sendError(w, "Already following this user", http.StatusConflict)
		return
	}

	if err := h.repo.Follow(followerID, followingID); err != nil {
		h.sendError(w, "Unable to follow user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "followed"})
	if h.publisher != nil {
		follower, followerErr := h.repo.GetUserByID(followerID)
		if followerErr == nil && follower != nil {
			h.publisher.Publish(realtime.Event{Type: "user.followed", Data: map[string]any{
				"follower_id":    followerID,
				"following_id":   followingID,
				"username":       follower.Username,
				"target_user_id": followingID,
			}})
		}
	}
}

func (h *FollowHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	followerID, ok := middleware.GetUserID(r)
	if !ok {
		h.sendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		h.sendError(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	followingID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		h.sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Unfollow(followerID, followingID); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "unfollowed"})
	if h.publisher != nil {
		h.publisher.Publish(realtime.Event{Type: "user.unfollowed", Data: map[string]int{"follower_id": followerID, "following_id": followingID}})
	}
}

func (h *FollowHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		h.sendError(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		h.sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	following, err := h.repo.GetFollowing(userID)
	if err != nil {
		h.sendError(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var users []map[string]interface{}
	for _, id := range following {
		user, _ := h.repo.GetUserByID(id)
		if user != nil {
			users = append(users, map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"name":     user.Name,
				"bio":      user.Bio,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *FollowHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		h.sendError(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		h.sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	followers, err := h.repo.GetFollowers(userID)
	if err != nil {
		h.sendError(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var users []map[string]interface{}
	for _, id := range followers {
		user, _ := h.repo.GetUserByID(id)
		if user != nil {
			users = append(users, map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"name":     user.Name,
				"bio":      user.Bio,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *FollowHandler) sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
