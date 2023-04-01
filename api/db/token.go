package db

import (
	"time"
)

type Token struct {
	ID          string
	ExpiresAt   time.Time
	ClientID    string
	UserAccount UserAccount
}

func (t *Token) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":           t.ID,
		"client_id":    t.ClientID,
		"user_account": t.UserAccount,
		"expires_at":   t.ExpiresAt,
	}
}
