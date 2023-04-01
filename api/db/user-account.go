package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/tashima42/go-oauth2-server/helpers"
)

type UserAccount struct {
	ID       string         `db:"id"`
	Username string         `db:"username"`
	Password string         `db:"password"`
	Scopes   pq.StringArray `db:"scopes"`
}

func UserAccountFromMap(m map[string]interface{}) UserAccount {
	log.Println(m)
	return UserAccount{
		ID:       m["id"].(string),
		Username: m["username"].(string),
	}
}

func (u *UserAccount) ScopesToSlice() []helpers.Scope {
	var scopes []helpers.Scope
	for _, s := range u.Scopes {
		scopes = append(scopes, helpers.Scope(s))
	}
	return scopes
}

func (r *Repo) CreateUserAccountTxx(tx *sqlx.Tx, u UserAccount) error {
	query := "INSERT INTO user_accounts(username, password, scopes) VALUES($1, $2, $3);"
	_, err := tx.Exec(query, u.Username, u.Password, u.Scopes)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetUserAccountByUsernameTxx(tx *sqlx.Tx, username string) (*UserAccount, error) {
	var u UserAccount
	query := "SELECT id, username, password, scopes FROM user_accounts WHERE username=$1 LIMIT 1;"
	err := tx.Get(&u, query, username)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repo) GetUserAccountByIDTxx(tx *sqlx.Tx, ID string) (*UserAccount, error) {
	var u UserAccount
	query := "SELECT id, username, password, scopes FROM user_accounts WHERE id=$1 LIMIT 1;"
	err := tx.Get(&u, query, ID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
