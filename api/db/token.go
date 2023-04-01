package db

import (
	"time"

	"github.com/tashima42/go-oauth2-server/helpers"
)

type Token struct {
	ExpiresAt   time.Time
	ClientID    string
	UserAccount UserAccount
	Scopes      []helpers.Scope
}

func (t *Token) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"clientID": t.ClientID,
		"userAccount": map[string]interface{}{
			"username": t.UserAccount.Username,
			"id":       t.UserAccount.ID,
		},
		"scopes": t.Scopes,
		"exp":    t.ExpiresAt.Unix(),
	}
}
