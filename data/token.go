package data

import (
	"database/sql"
	"time"

	"github.com/tashima42/go-oauth2-server/helpers"
)

type Token struct {
	ID                    int    `json:"id"`
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	Active                bool   `json:"active"`
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	ClientId              int
	UserAccountId         int
}

func (t *Token) CreateToken(db *sql.DB) error {
	t.AccessToken = helpers.GenerateSecureToken(64)
	t.RefreshToken = helpers.GenerateSecureToken(64)
	t.AccessTokenExpiresAt = helpers.NowPlusSeconds(86400)
	t.RefreshTokenExpiresAt = helpers.NowPlusSeconds(2628288)

	return db.QueryRow(
		"INSERT INTO tokens(access_token, access_token_expires_at, refresh_token, refresh_token_expires_at, client_id, user_account_id) VALUES($1, $2, $3, $4, $5, $6) RETURNING id;",
		t.AccessToken,
		helpers.FormatDateIso(t.AccessTokenExpiresAt),
		t.RefreshToken,
		helpers.FormatDateIso(t.RefreshTokenExpiresAt),
		t.ClientId,
		t.UserAccountId,
	).Scan(&t.ID)
}

func (t *Token) GetByRefreshToken(db *sql.DB) error {
	var accessTokenExpiresAt string
	var refreshTokenExpiresAt string
	err := db.QueryRow(
		"SELECT id, access_token, access_token_expires_at, refresh_token, refresh_token_expires_at, client_id, user_account_id, active FROM tokens WHERE refresh_token=$1 LIMIT 1;",
		t.RefreshToken,
	).Scan(&t.ID, &t.AccessToken, &accessTokenExpiresAt, &t.RefreshToken, &refreshTokenExpiresAt, &t.ClientId, &t.UserAccountId, &t.Active)
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
	).Scan(&t.ID, &t.AccessToken, &accessTokenExpiresAt, &t.RefreshToken, &refreshTokenExpiresAt, &t.ClientId, &t.UserAccountId, &t.Active)
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
