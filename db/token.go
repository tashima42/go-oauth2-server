package db

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tashima42/go-oauth2-server/helpers"
)

type Token struct {
	ID          int         `json:"id"`
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

func (r *Repo) CreateTokenTxx(tx *sqlx.Tx, token Token) error {
	query := "INSERT INTO tokens(access_token, access_token_expires_at, refresh_token, refresh_token_expires_at, client_id, user_account_id) VALUES($1, $2, $3, $4, $5, $6);"
	_, err := tx.Exec(query, token.AccessToken, token.AccessTokenExpiresAt, token.RefreshToken, token.RefreshTokenExpiresAt, token.ClientID, token.UserAccountID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetByRefreshToken(db *sql.DB) error {
	var accessTokenExpiresAt string
	var refreshTokenExpiresAt string
	err := db.QueryRow(
		"SELECT id, access_token, access_token_expires_at, refresh_token, refresh_token_expires_at, client_id, user_account_id, active FROM tokens WHERE refresh_token=$1 LIMIT 1;",
		t.RefreshToken,
	).Scan(&t.ID, &t.AccessToken, &accessTokenExpiresAt, &t.RefreshToken, &refreshTokenExpiresAt, &t.ClientID, &t.UserAccountID, &t.Active)
	if err != nil {
		return err
	}
	t.AccessTokenExpiresAt, err = helpers.ParseDateIso(accessTokenExpiresAt)
	if err != nil {
		return err
	}
	t.RefreshTokenExpiresAt, err = helpers.ParseDateIso(refreshTokenExpiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (t *Token) GetByAccessToken(db *sql.DB) error {
	var accessTokenExpiresAt string
	var refreshTokenExpiresAt string
	err := db.QueryRow(
		"SELECT id, access_token, access_token_expires_at, refresh_token, refresh_token_expires_at, client_id, user_account_id, active FROM tokens WHERE access_token=$1 LIMIT 1;",
		t.AccessToken,
	).Scan(&t.ID, &t.AccessToken, &accessTokenExpiresAt, &t.RefreshToken, &refreshTokenExpiresAt, &t.ClientID, &t.UserAccountID, &t.Active)
	if err != nil {
		return err
	}
	t.AccessTokenExpiresAt, err = helpers.ParseDateIso(accessTokenExpiresAt)
	if err != nil {
		return err
	}
	t.RefreshTokenExpiresAt, err = helpers.ParseDateIso(refreshTokenExpiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (t *Token) Disable(db *sql.DB) error {
	_, err := db.Exec("UPDATE tokens SET active=false WHERE id=$1;", t.ID)
	return err
}
