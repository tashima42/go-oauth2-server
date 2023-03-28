package db

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tashima42/go-oauth2-server/helpers"
)

type AuthorizationCode struct {
	ID            string    `json:"id" db:"id"`
	Code          string    `json:"code" db:"code"`
	Active        bool      `json:"active" db:"active"`
	ExpiresAt     time.Time `json:"expires_at" db:"expires_at"`
	RedirectURI   string    `json:"redirect_uri" db:"redirect_uri"`
	ClientID      string    `db:"client_id"`
	UserAccountID string    `db:"user_account_id"`
}

func (r *Repo) CreateAuthorizationCodeTxx(tx *sqlx.Tx, ac AuthorizationCode) error {
	query := "INSERT INTO authorization_codes(code, expires_at, redirect_uri, client_id, user_account_id) VALUES($1, $2, $3, $4, $5);"
	_, err := tx.Exec(query, ac.Code, helpers.FormatDateIso(ac.ExpiresAt), ac.RedirectURI, ac.ClientID, ac.UserAccountID)
	if err != nil {
		return err
	}
	return nil
}

func (ac *AuthorizationCode) GetByCode(db *sql.DB) error {
	var expiresAt string
	err := db.QueryRow(
		"SELECT id, code, expires_at, redirect_uri,client_id, user_account_id, active FROM authorization_codes WHERE code=$1 LIMIT 1;",
		ac.Code,
	).Scan(&ac.ID, &ac.Code, &expiresAt, &ac.RedirectURI, &ac.ClientID, &ac.UserAccountID, &ac.Active)
	if err != nil {
		return err
	}
	ac.ExpiresAt, err = helpers.ParseDateIso(expiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (ac *AuthorizationCode) Disable(db *sql.DB) error {
	_, err := db.Exec("UPDATE authorization_codes SET active=false WHERE id=$1;", ac.ID)
	return err
}
