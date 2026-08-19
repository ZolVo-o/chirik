package handlers

import (
	"chirik/internal/middleware"
	"chirik/internal/models"
	"chirik/internal/realtime"
	"chirik/internal/repository"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type MessageHandler struct {
	repo      *repository.Repository
	publisher realtime.Publisher
}

func NewMessageHandler(repo *repository.Repository, publisher realtime.Publisher) *MessageHandler {
	return &MessageHandler{repo: repo, publisher: publisher}
}

type conversationRequest struct {
	UserID int `json:"user_id"`
}

type messageRequest struct {
	Content string `json:"content"`
}

func (h *MessageHandler) Conversations(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodGet {
		conversations, err := h.repo.GetConversations(userID)
		if err != nil {
			sendJSONError(w, "Internal error", http.StatusInternalServerError)
			return
		}
		sendJSON(w, conversations)
		return
	}
	if r.Method != http.MethodPost {
		sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req conversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID <= 0 || req.UserID == userID {
		sendJSONError(w, "Invalid user", http.StatusBadRequest)
		return
	}
	otherUser, err := h.repo.GetUserByID(req.UserID)
	if err != nil {
		sendJSONError(w, "Internal error", http.StatusInternalServerError)
		return
	}
	if otherUser == nil {
		sendJSONError(w, "User not found", http.StatusNotFound)
		return
	}
	conversation, err := h.repo.CreateConversation(userID, req.UserID)
	if err != nil {
		sendJSONError(w, "Unable to create conversation", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	sendJSON(w, conversation)
}

func (h *MessageHandler) ConversationMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	conversationID, ok := pathID(r.URL.Path, "/api/conversations/", "/messages")
	if !ok {
		sendJSONError(w, "Invalid conversation", http.StatusBadRequest)
		return
	}
	isMember, err := h.repo.IsConversationMember(conversationID, userID)
	if err != nil {
		sendJSONError(w, "Internal error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		sendJSONError(w, "Conversation not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		messages, err := h.repo.GetMessages(conversationID)
		if err != nil {
			sendJSONError(w, "Internal error", http.StatusInternalServerError)
			return
		}
		sendJSON(w, messages)
		return
	}
	if r.Method != http.MethodPost {
		sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req messageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	req.Content = strings.TrimSpace(req.Content)
	if len(req.Content) == 0 || len(req.Content) > 2000 {
		sendJSONError(w, "Message must contain 1-2000 characters", http.StatusBadRequest)
		return
	}
	user, err := h.repo.GetUserByID(userID)
	if err != nil || user == nil {
		sendJSONError(w, "User not found", http.StatusNotFound)
		return
	}
	message := &models.Message{ConversationID: conversationID, SenderID: userID, SenderUsername: user.Username, Content: req.Content}
	if err := h.repo.CreateMessage(message); err != nil {
		sendJSONError(w, "Unable to send message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	sendJSON(w, message)
	h.publishMessage("message.created", message, userID)
}

func (h *MessageHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	h.mutateMessage(w, r, false)
}

func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	h.mutateMessage(w, r, true)
}

func (h *MessageHandler) mutateMessage(w http.ResponseWriter, r *http.Request, remove bool) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	messageID, ok := trailingID(r.URL.Path)
	if !ok {
		sendJSONError(w, "Invalid message", http.StatusBadRequest)
		return
	}
	var message *models.Message
	var err error
	if remove {
		if r.Method != http.MethodDelete {
			sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		message, err = h.repo.DeleteMessage(messageID, userID)
	} else {
		if r.Method != http.MethodPut {
			sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req messageRequest
		if decodeErr := json.NewDecoder(r.Body).Decode(&req); decodeErr != nil {
			sendJSONError(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		req.Content = strings.TrimSpace(req.Content)
		if len(req.Content) == 0 || len(req.Content) > 2000 {
			sendJSONError(w, "Message must contain 1-2000 characters", http.StatusBadRequest)
			return
		}
		message, err = h.repo.UpdateMessage(messageID, userID, req.Content)
	}
	if errors.Is(err, sql.ErrNoRows) {
		sendJSONError(w, "Message not found or not owned by you", http.StatusNotFound)
		return
	}
	if err != nil {
		sendJSONError(w, "Internal error", http.StatusInternalServerError)
		return
	}
	sendJSON(w, message)
	h.publishMessage("message.updated", message, userID)
}

func (h *MessageHandler) publish(eventType string, data any) {
	if h.publisher != nil {
		h.publisher.Publish(realtime.Event{Type: eventType, Data: data})
	}
}

func (h *MessageHandler) publishMessage(eventType string, message *models.Message, senderID int) {
	// SSE-поток общий, поэтому не отправляем через него текст личного сообщения.
	data := map[string]any{
		"message_id":      message.ID,
		"conversation_id": message.ConversationID,
		"sender_id":       message.SenderID,
	}
	conversation, err := h.repo.GetConversation(message.ConversationID, senderID)
	if err == nil && conversation != nil && conversation.OtherUser != nil {
		data["target_user_id"] = conversation.OtherUser.ID
	}
	h.publish(eventType, data)
}

func pathID(path, prefix, suffix string) (int, bool) {
	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, suffix) {
		return 0, false
	}
	value := strings.TrimSuffix(strings.TrimPrefix(path, prefix), suffix)
	id, err := strconv.Atoi(value)
	return id, err == nil && id > 0
}

func trailingID(path string) (int, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return 0, false
	}
	id, err := strconv.Atoi(parts[len(parts)-1])
	return id, err == nil && id > 0
}

func sendJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}

func sendJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	sendJSON(w, map[string]string{"error": message})
}
