package db

import (
	"time"
)

type Token struct {
	ExpiresAt   time.Time
	ClientID    string
	UserAccount UserAccount
}

func (t *Token) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"clientID": t.ClientID,
		"userAccount": map[string]interface{}{
			"username": t.UserAccount.Username,
			"id":       t.UserAccount.ID,
		},
		"exp": t.ExpiresAt.Unix(),
	}
}
