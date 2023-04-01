package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type AuthorizationCode struct {
	ID            string    `db:"id"`
	Code          string    `db:"code"`
	Active        bool      `db:"active"`
	ExpiresAt     time.Time `db:"expires_at"`
	RedirectURI   string    `db:"redirect_uri"`
	ClientID      string    `db:"client_id"`
	UserAccountID string    `db:"user_account_id"`
}

func (r *Repo) CreateAuthorizationCodeTxx(tx *sqlx.Tx, ac AuthorizationCode) error {
	query := "INSERT INTO authorization_codes(code, expires_at, redirect_uri, client_id, user_account_id) VALUES($1, $2, $3, $4, $5);"
	_, err := tx.Exec(query, ac.Code, ac.ExpiresAt, ac.RedirectURI, ac.ClientID, ac.UserAccountID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetAuthorizationCodeByCodeTxx(tx *sqlx.Tx, code string) (*AuthorizationCode, error) {
	var ac AuthorizationCode
	query := "SELECT id, code, expires_at, redirect_uri,client_id, user_account_id, active FROM authorization_codes WHERE code=$1 LIMIT 1;"
	err := tx.Get(&ac, query, code)
	if err != nil {
		return nil, err
	}
	return &ac, nil
}

func (r *Repo) DisableAuthorizationCodeByIDTxx(tx *sqlx.Tx, ID string) error {
	query := "UPDATE authorization_codes SET active=false WHERE id=$1;"
	_, err := tx.Exec(query, ID)
	return err
}
