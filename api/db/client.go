package db

import (
	"github.com/jmoiron/sqlx"
)

type Client struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	ClientID     string `json:"client_id" db:"client_id"`
	ClientSecret string `json:"client_secret" db:"client_secret"`
	RedirectURI  string `json:"redirect_uri" db:"redirect_uri"`
}

func (r *Repo) CreateClientTxx(tx *sqlx.Tx, c Client) error {
	query := "INSERT INTO clients(name, client_id, client_secret, redirect_uri) VALUES($1, $2, $3, $4);"
	_, err := tx.Exec(query, c.Name, c.ClientID, c.ClientSecret, c.RedirectURI)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetClientByClientIDTxx(tx *sqlx.Tx, clientID string) (*Client, error) {
	var c Client
	query := "SELECT id, name, client_id, client_secret, redirect_uri FROM clients WHERE client_id=$1 LIMIT 1;"
	err := tx.Get(&c, query, clientID)
	if err != nil {
		return nil, err
	}
	return &c, err
}
