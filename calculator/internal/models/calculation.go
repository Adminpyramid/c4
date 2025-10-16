package models

import "time"

type Calculation struct {
	ID         string    `json:"id"`
	SessionID  string    `json:"session_id"`
	UserID     string    `json:"user_id"`
	Expression string    `json:"expression"`
	Result     string    `json:"result"`
	CreatedAt  time.Time `json:"created_at"`
}
