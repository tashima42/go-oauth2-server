package db

import (
	"time"
)

type Token struct {
	ID          string      `json:"id"`
	ExpiresAt   time.Time   `json:"expires_at"`
	ClientID    string      `json:"client_id"`
	UserAccount UserAccount `json:"user_account"`
}

func (t *Token) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":           t.ID,
		"client_id":    t.ClientID,
		"user_account": t.UserAccount,
		"expires_at":   t.ExpiresAt,
	}
}
