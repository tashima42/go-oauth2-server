package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type UserAccount struct {
	ID           string `db:"id"`
	Username     string `db:"username"`
	Password     string `db:"password"`
	Country      string `db:"country"`
	SubscriberID string `db:"subscriber_id"`
}

func (r *Repo) CreateUserAccountTxx(tx *sql.Tx, u UserAccount) error {
	query := "INSERT INTO user_accounts(username, password, country, subscriber_id) VALUES($1, $2, $3, $4) RETURNING id;"
	_, err := tx.Exec(query, u.Username, u.Password, u.Country, u.SubscriberID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetUserAccountByUsernameAndCountryTxx(tx *sqlx.Tx, username string, country string) (*UserAccount, error) {
	var u UserAccount
	query := "SELECT id, username, password, country, subscriber_id FROM user_accounts WHERE username=$1 AND country=$2 LIMIT 1;"
	err := tx.Get(&u, query, username, country)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repo) GetUserAccountByIDTxx(tx *sqlx.Tx, ID string) (*UserAccount, error) {
	var u UserAccount
	query := "SELECT id, username, password, country, subscriber_id FROM user_accounts WHERE id=$1 LIMIT 1;"
	err := tx.Get(&u, query, ID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
