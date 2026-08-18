package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "chirik/internal/models"
    "chirik/internal/repository"
    "chirik/internal/middleware"
)

type TweetHandler struct {
    repo *repository.Repository
}

func NewTweetHandler(repo *repository.Repository) *TweetHandler {
    return &TweetHandler{repo: repo}
}

type CreateTweetRequest struct {
    Content string `json:"content"`
}

func (h *TweetHandler) CreateTweet(w http.ResponseWriter, r *http.Request) {
    // Проверяем только метод
    if r.Method != http.MethodPost {
        h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Получаем user_id из контекста (уже сохранён middleware)
    userID, ok := middleware.GetUserID(r)
    if !ok || userID == 0 {
        h.sendError(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var req CreateTweetRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.sendError(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Валидация
    if len(req.Content) < 3 {
        h.sendError(w, "Tweet must be at least 3 characters", http.StatusBadRequest)
        return
    }
    if len(req.Content) > 280 {
        h.sendError(w, "Tweet must be less than 280 characters", http.StatusBadRequest)
        return
    }

    // Получаем пользователя
    user, err := h.repo.GetUserByID(userID)
    if err != nil || user == nil {
        h.sendError(w, "User not found", http.StatusNotFound)
        return
    }

    // Создаём твит
    tweet := &models.Tweet{
        UserID:   userID,
        Username: user.Username,
        Content:  req.Content,
    }

    if err := h.repo.CreateTweet(tweet); err != nil {
        h.sendError(w, "Internal error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(tweet)
}

func (h *TweetHandler) GetAllTweets(w http.ResponseWriter, r *http.Request) {
    tweets, err := h.repo.GetAllTweets()
    if err != nil {
        h.sendError(w, "Internal error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tweets)
}

func (h *TweetHandler) GetTweetsByUser(w http.ResponseWriter, r *http.Request) {
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

    tweets, err := h.repo.GetTweetsByUser(userID)
    if err != nil {
        h.sendError(w, "Internal error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tweets)
}

func (h *TweetHandler) LikeTweet(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 4 {
        h.sendError(w, "Invalid URL", http.StatusBadRequest)
        return
    }

    tweetID, err := strconv.Atoi(parts[len(parts)-1])
    if err != nil {
        h.sendError(w, "Invalid tweet ID", http.StatusBadRequest)
        return
    }

    if err := h.repo.LikeTweet(tweetID); err != nil {
        h.sendError(w, "Tweet not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "liked"})
}

func (h *TweetHandler) sendError(w http.ResponseWriter, message string, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}
