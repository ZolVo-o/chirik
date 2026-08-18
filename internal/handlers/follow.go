package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "chirik/internal/repository"
)

type FollowHandler struct {
    repo *repository.Repository
}

func NewFollowHandler(repo *repository.Repository) *FollowHandler {
    return &FollowHandler{repo: repo}
}

func (h *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
    followerID := r.Context().Value("user_id").(int)

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

    if err := h.repo.Follow(followerID, followingID); err != nil {
        h.sendError(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "followed"})
}

func (h *FollowHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
    followerID := r.Context().Value("user_id").(int)

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
}

func (h *FollowHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
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
