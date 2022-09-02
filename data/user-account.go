package data

import (
	"database/sql"
)

type UserAccount struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Country      string `json:"country"`
	SubscriberId string `json:"subscriber_id"`
}

func (u *UserAccount) CreateUserAccount(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO user_accounts(username, password, country, subscriber_id) VALUES($1, $2, $3, $4) RETURNING id;",
		u.Username,
		u.Password,
		u.Country,
		u.SubscriberId,
	).Scan(&u.ID)
}

func (u *UserAccount) GetByUsernameAndCountry(db *sql.DB) error {
	return db.QueryRow(
		"SELECT id, username, password, country, subscriber_id FROM user_accounts WHERE username=$1 AND country=$2 LIMIT 1;",
		u.Username, u.Country,
	).Scan(&u.ID, &u.Username, &u.Password, &u.Country, &u.SubscriberId)
}

func (u *UserAccount) GetById(db *sql.DB) error {
	return db.QueryRow(
		"SELECT id, username, password, country, subscriber_id FROM user_accounts WHERE id=$1 LIMIT 1;",
		u.ID,
	).Scan(&u.ID, &u.Username, &u.Password, &u.Country, &u.SubscriberId)
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
