package models

import "time"

type Follow struct {
    ID          int       `json:"id"`
    FollowerID  int       `json:"follower_id"`
    FollowingID int       `json:"following_id"`
    CreatedAt   time.Time `json:"created_at"`
}
