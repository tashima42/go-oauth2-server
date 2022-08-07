package data

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type UserAccount struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Country      string `json:"country"`
	SubscriberId string `json:"subscriber_id"`
}

func (u *UserAccount) CreateUserAccount(db *sql.DB) error {
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}

	u.Password = string(password)
	return db.QueryRow(
		"INSERT INTO user_accounts(username, password, country, subscriber_id) VALUES($1, $2, $3, $4) RETURNING id;",
		u.Username,
		u.Password,
		u.Country,
		u.SubscriberId,
	).Scan(&u.ID)
}

func (u *UserAccount) GetByUsernameAndCountry(db *sql.DB, username string, country string) error {
	return db.QueryRow(
		"SELECT id, username, password, country, subscriber_id FROM user_accounts WHERE username=$1 AND country=$2 LIMIT 1;",
		username, country,
	).Scan(&u.ID, &u.Username, &u.Password, &u.Country, &u.SubscriberId)
}

func (u *UserAccount) GetById(db *sql.DB, id int) error {
	return db.QueryRow(
		"SELECT id, username, password, country, subscriber_id FROM user_accounts WHERE id=$1 LIMIT 1;",
		id,
	).Scan(&u.ID, &u.Username, &u.Password, &u.Country, &u.SubscriberId)
}
