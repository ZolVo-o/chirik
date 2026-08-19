package handlers

import (
	"chirik/internal/middleware"
	"chirik/internal/models"
	"chirik/internal/realtime"
	"chirik/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type TweetHandler struct {
	repo      *repository.Repository
	publisher realtime.Publisher
}

func NewTweetHandler(repo *repository.Repository, publisher realtime.Publisher) *TweetHandler {
	return &TweetHandler{repo: repo, publisher: publisher}
}

type CreateTweetRequest struct {
	Content string `json:"content"`
}

func (h *TweetHandler) CreateTweet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	req.Content = strings.TrimSpace(req.Content)

	if len(req.Content) < 3 {
		h.sendError(w, "Tweet must be at least 3 characters", http.StatusBadRequest)
		return
	}
	if len(req.Content) > 280 {
		h.sendError(w, "Tweet must be less than 280 characters", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(userID)
	if err != nil || user == nil {
		h.sendError(w, "User not found", http.StatusNotFound)
		return
	}

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
	if h.publisher != nil {
		h.publisher.Publish(realtime.Event{Type: "tweet.created", Data: tweet})
	}
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

	tweets, err := h.repo.GetTweetsByUser(userID)
	if err != nil {
		h.sendError(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tweets)
}

func (h *TweetHandler) LikeTweet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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

	userID, ok := middleware.GetUserID(r)
	if !ok {
		h.sendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	liked, err := h.repo.LikeTweet(tweetID, userID)
	if err != nil {
		h.sendError(w, "Tweet not found", http.StatusNotFound)
		return
	}
	if !liked {
		h.sendError(w, "Tweet already liked", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "liked"})
	if h.publisher != nil {
		tweet, tweetErr := h.repo.GetTweetByID(tweetID)
		liker, likerErr := h.repo.GetUserByID(userID)
		if tweetErr == nil && likerErr == nil && tweet != nil && liker != nil {
			h.publisher.Publish(realtime.Event{Type: "tweet.liked", Data: map[string]any{
				"tweet_id":       tweetID,
				"user_id":        userID,
				"username":       liker.Username,
				"target_user_id": tweet.UserID,
			}})
		}
	}
}

func (h *TweetHandler) ViewTweet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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

	userID, ok := middleware.GetUserID(r)
	if !ok {
		h.sendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	viewed, err := h.repo.ViewTweet(tweetID, userID)
	if err != nil {
		h.sendError(w, "Tweet not found", http.StatusNotFound)
		return
	}
	if !viewed {
		tweet, err := h.repo.GetTweetByID(tweetID)
		if err != nil || tweet == nil {
			h.sendError(w, "Tweet not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tweet)
		if h.publisher != nil {
			h.publisher.Publish(realtime.Event{Type: "tweet.viewed", Data: map[string]int{"tweet_id": tweetID}})
		}
		return
	}

	tweet, err := h.repo.GetTweetByID(tweetID)
	if err != nil {
		h.sendError(w, "Tweet not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tweet)
	if viewed && h.publisher != nil {
		h.publisher.Publish(realtime.Event{Type: "tweet.viewed", Data: map[string]int{"tweet_id": tweetID}})
	}
}

func (h *TweetHandler) sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
