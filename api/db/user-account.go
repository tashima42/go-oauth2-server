package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AccountType string

const (
	UserAccountType AccountType = "user"
	DevAccountType  AccountType = "dev"
)

type UserAccount struct {
	ID       string         `db:"id"`
	Username string         `db:"username"`
	Password string         `db:"password"`
	Type     AccountType    `db:"type"`
	Scopes   pq.StringArray `db:"scopes"`
}

func UserAccountFromMap(m map[string]interface{}) UserAccount {
	userAccountType := m["type"].(string)
	return UserAccount{
		ID:       m["id"].(string),
		Username: m["username"].(string),
		Type:     AccountType(userAccountType),
	}
}

func (u *UserAccount) ScopesToSlice() []string {
	var scopes []string
	for _, s := range u.Scopes {
		scopes = append(scopes, s)
	}
	return scopes
}

func (r *Repo) CreateUserAccountTxx(tx *sqlx.Tx, u UserAccount) error {
	query := "INSERT INTO user_accounts(username, password, scopes, type) VALUES($1, $2, $3, $4);"
	_, err := tx.Exec(query, u.Username, u.Password, u.Scopes, u.Type)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetUserAccountByUsernameTxx(tx *sqlx.Tx, username string) (*UserAccount, error) {
	var u UserAccount
	query := "SELECT id, username, password, scopes, type FROM user_accounts WHERE username=$1 LIMIT 1;"
	err := tx.Get(&u, query, username)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repo) GetUserAccountByIDTxx(tx *sqlx.Tx, ID string) (*UserAccount, error) {
	var u UserAccount
	query := "SELECT id, username, password, scopes, type FROM user_accounts WHERE id=$1 LIMIT 1;"
	err := tx.Get(&u, query, ID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
