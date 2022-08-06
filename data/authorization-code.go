package data

import (
	"database/sql"
	"time"

	"github.com/tashima42/go-oauth2-server/helpers"
)

type AuthorizationCode struct {
	ID            int    `json:"id"`
	Code          string `json:"code"`
	Active        bool   `json:"active"`
	expiresAt     time.Time
	RedirectUri   string
	ClientId      int
	UserAccountId int
}

func (ac *AuthorizationCode) CreateAuthorizationCode(db *sql.DB) error {
	ac.Code = helpers.RandStringBytes(128)
	ac.expiresAt = helpers.NowPlusSeconds(86400)

	return db.QueryRow(
		"INSERT INTO authorization_codes(code, expires_at, redirect_uri, client_id, user_account_id) VALUES($1, $2, $3, $4, $5) RETURNING id;",
		ac.Code,
		helpers.FormatDateIso(ac.expiresAt),
		ac.RedirectUri,
		ac.ClientId,
		ac.UserAccountId,
	).Scan(&ac.ID)
}

func (ac *AuthorizationCode) GetByCode(db *sql.DB, code string) error {
	var expiresAt string
	err := db.QueryRow(
		"SELECT id, code, expires_at, redirect_uri,client_id, user_account_id, active FROM authorization_codes WHERE code=$1 LIMIT 1;",
		code,
	).Scan(&ac.ID, &ac.Code, &expiresAt, &ac.RedirectUri, &ac.ClientId, &ac.UserAccountId, &ac.Active)
	if err != nil {
		return err
	}
	return nil
}

func (ac *AuthorizationCode) Disable(db *sql.DB) error {
	_, err := db.Exec("UPDATE authorization_codes SET active=false WHERE id=$1;", ac.ID)
	return err
}
