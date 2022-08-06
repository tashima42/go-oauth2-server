package data

import (
	"database/sql"
)

type Client struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
}

func (c *Client) CreateClient(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO clients(name, client_id, client_secret, redirect_uri) VALUES($1, $2, $3, $4) RETURNING id;",
		c.Name,
		c.ClientId,
		c.ClientSecret,
		c.RedirectUri,
	).Scan(&c.ID)
}

func (c *Client) GetByClientId(db *sql.DB, clientId string) error {
	return db.QueryRow(
		"SELECT id, name, client_id, client_secret, redirect_uri FROM clients WHERE client_id=$1 LIMIT 1;",
		clientId,
	).Scan(&c.ID, &c.Name, &c.ClientId, &c.ClientSecret, &c.RedirectUri)
}
