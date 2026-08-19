package models

import "time"

type Conversation struct {
	ID          int       `json:"id"`
	OtherUser   *User     `json:"other_user"`
	LastMessage *Message  `json:"last_message,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Message struct {
	ID             int       `json:"id"`
	ConversationID int       `json:"conversation_id"`
	SenderID       int       `json:"sender_id"`
	SenderUsername string    `json:"sender_username"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Deleted        bool      `json:"deleted"`
}
