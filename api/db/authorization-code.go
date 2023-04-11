package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AuthorizationCode struct {
	ID            string      `db:"id"`
	Code          string      `db:"code"`
	Active        bool        `db:"active"`
	ExpiresAt     pq.NullTime `db:"expires_at"`
	RedirectURI   string      `db:"redirect_uri"`
	ClientID      string      `db:"client_id"`
	UserAccountID string      `db:"user_account_id"`
}

func (r *Repo) CreateAuthorizationCodeTxx(tx *sqlx.Tx, ac AuthorizationCode) error {
	query := "INSERT INTO authorization_codes(code, expires_at, redirect_uri, client_id, user_account_id) VALUES($1, $2, $3, $4, $5);"
	_, err := tx.Exec(query, ac.Code, ac.ExpiresAt, ac.RedirectURI, ac.ClientID, ac.UserAccountID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetAuthorizationCodeByCodeAndClientTxx(tx *sqlx.Tx, code string, clientID string) (*AuthorizationCode, error) {
	var ac AuthorizationCode
	query := "SELECT id, code, expires_at, redirect_uri,client_id, user_account_id, active FROM authorization_codes WHERE code=$1 AND client_id=$2 LIMIT 1;"
	err := tx.Get(&ac, query, code, clientID)
	if err != nil {
		log.Println("NEW ERROR", err)
		return nil, err
	}
	return &ac, nil
}

func (r *Repo) DisableAuthorizationCodeByIDTxx(tx *sqlx.Tx, ID string) error {
	query := "UPDATE authorization_codes SET active=false WHERE id=$1;"
	_, err := tx.Exec(query, ID)
	return err
}
