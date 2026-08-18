package handlers

import (
    "encoding/json"
    "net/http"
    "golang.org/x/crypto/bcrypt"
    "chirik/internal/auth"
    "chirik/internal/models"
    "chirik/internal/repository"
)

type AuthHandler struct {
    repo *repository.Repository
}

func NewAuthHandler(repo *repository.Repository) *AuthHandler {
    return &AuthHandler{repo: repo}
}

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Name     string `json:"name"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type UpdateProfileRequest struct {
    Name string `json:"name"`
    Bio  string `json:"bio"`
}

type AuthResponse struct {
    Token string      `json:"token"`
    User  models.User `json:"user"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.sendError(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if req.Username == "" || len(req.Username) < 3 {
        h.sendError(w, "Username must be at least 3 characters", http.StatusBadRequest)
        return
    }
    if req.Email == "" {
        h.sendError(w, "Email required", http.StatusBadRequest)
        return
    }
    if len(req.Password) < 6 {
        h.sendError(w, "Password must be at least 6 characters", http.StatusBadRequest)
        return
    }
    if req.Name == "" {
        h.sendError(w, "Name required", http.StatusBadRequest)
        return
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        h.sendError(w, "Internal error", http.StatusInternalServerError)
        return
    }

    user := &models.User{
        Username: req.Username,
        Email:    req.Email,
        Password: string(hashed),
        Name:     req.Name,
        Bio:      "",
    }

    if err := h.repo.CreateUser(user); err != nil {
        h.sendError(w, err.Error(), http.StatusConflict)
        return
    }

    token, err := auth.GenerateToken(user.ID)
    if err != nil {
        h.sendError(w, "Internal error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(AuthResponse{Token: token, User: *user})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.sendError(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    user, err := h.repo.GetUserByEmail(req.Email)
    if err != nil || user == nil {
        h.sendError(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        h.sendError(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    token, err := auth.GenerateToken(user.ID)
    if err != nil {
        h.sendError(w, "Internal error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(AuthResponse{Token: token, User: *user})
}

func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("user_id").(int)
    if !ok || userID == 0 {
        h.sendError(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    user, err := h.repo.GetUserByID(userID)
    if err != nil || user == nil {
        h.sendError(w, "User not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("user_id").(int)
    if !ok || userID == 0 {
        h.sendError(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var req UpdateProfileRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.sendError(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    user, err := h.repo.GetUserByID(userID)
    if err != nil || user == nil {
        h.sendError(w, "User not found", http.StatusNotFound)
        return
    }

    if req.Name != "" {
        user.Name = req.Name
    }
    user.Bio = req.Bio

    if err := h.repo.UpdateUser(user); err != nil {
        h.sendError(w, "Internal error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) sendError(w http.ResponseWriter, message string, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}
