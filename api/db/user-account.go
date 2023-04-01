package db

import (
	"github.com/jmoiron/sqlx"
)

type UserAccount struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

func (r *Repo) CreateUserAccountTxx(tx *sqlx.Tx, u UserAccount) error {
	query := "INSERT INTO user_accounts(username, password) VALUES($1, $2);"
	_, err := tx.Exec(query, u.Username, u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetUserAccountByUsernameTxx(tx *sqlx.Tx, username string) (*UserAccount, error) {
	var u UserAccount
	query := "SELECT id, username, password FROM user_accounts WHERE username=$1 LIMIT 1;"
	err := tx.Get(&u, query, username)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repo) GetUserAccountByIDTxx(tx *sqlx.Tx, ID string) (*UserAccount, error) {
	var u UserAccount
	query := "SELECT id, username, password FROM user_accounts WHERE id=$1 LIMIT 1;"
	err := tx.Get(&u, query, ID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
