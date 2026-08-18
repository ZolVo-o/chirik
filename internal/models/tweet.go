package models

import "time"

type Tweet struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Username  string    `json:"username"`
    Content   string    `json:"content"`
    Likes     int       `json:"likes"`
    CreatedAt time.Time `json:"created_at"`
}
