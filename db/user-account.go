package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type UserAccount struct {
	ID           string `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	Password     string `json:"password" db:"password"`
	Country      string `json:"country" db:"country"`
	SubscriberId string `json:"subscriber_id" db:"subscriber_id"`
}

func (r *Repo) CreateUserAccountTxx(tx *sql.Tx, u UserAccount) error {
	query := "INSERT INTO user_accounts(username, password, country, subscriber_id) VALUES($1, $2, $3, $4) RETURNING id;"
	_, err := tx.Exec(query, u.Username, u.Password, u.Country, u.SubscriberId)
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

func (u *UserAccount) GetByUsername(db *sql.DB) error {
	return db.QueryRow(
		"SELECT id, username, password, country, subscriber_id FROM user_accounts WHERE username=$1 LIMIT 1;",
		u.Username,
	).Scan(&u.ID, &u.Username, &u.Password, &u.Country, &u.SubscriberId)
}

func (u *UserAccount) GetBySubscriberId(db *sql.DB) error {
	return db.QueryRow(
		"SELECT id, username, password, country, subscriber_id FROM user_accounts WHERE subscriber_id=$1 LIMIT 1;",
		u.SubscriberId,
	).Scan(&u.ID, &u.Username, &u.Password, &u.Country, &u.SubscriberId)
}
