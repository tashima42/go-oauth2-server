package main

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

type Client struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
}

type AuthorizationCode struct {
	ID            int    `json:"id"`
	Code          string `json:"code"`
	Active        bool   `json:"active"`
	expiresAt     time.Time
	RedirectUri   string
	ClientId      int
	UserAccountId int
}

type UserAccount struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Country      string `json:"country"`
	SubscriberId string `json:"subscriber_id"`
}

type Token struct {
	ID                    int    `json:"id"`
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	Active                bool   `json:"active"`
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	ClientId              int
	UserAccountId         int
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

func (ac *AuthorizationCode) CreateAuthorizationCode(db *sql.DB) error {
	ac.Code = RandStringBytes(128)
	ac.expiresAt = nowPlusSeconds(86400)

	return db.QueryRow(
		"INSERT INTO authorization_codes(code, expires_at, redirect_uri, client_id, user_account_id) VALUES($1, $2, $3, $4, $5) RETURNING id;",
		ac.Code,
		formatDateIso(ac.expiresAt),
		ac.RedirectUri,
		ac.ClientId,
		ac.UserAccountId,
	).Scan(&ac.ID)
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

func (t *Token) CreateToken(db *sql.DB) error {
	t.AccessToken = RandStringBytes(128)
	t.RefreshToken = RandStringBytes(128)
	t.AccessTokenExpiresAt = nowPlusSeconds(86400)
	t.RefreshTokenExpiresAt = nowPlusSeconds(2628288)

	return db.QueryRow(
		"INSERT INTO tokens(access_token, access_token_expires_at, refresh_token, refresh_token_expires_at, client_id, user_account_id) VALUES($1, $2, $3, $4, $5, $6) RETURNING id;",
		t.AccessToken,
		formatDateIso(t.AccessTokenExpiresAt),
		t.RefreshToken,
		formatDateIso(t.RefreshTokenExpiresAt),
		t.ClientId,
		t.UserAccountId,
	).Scan(&t.ID)
}

func (c *Client) GetByClientId(db *sql.DB, clientId string) error {
	return db.QueryRow(
		"SELECT id, name, client_id, client_secret, redirect_uri FROM clients WHERE client_id=$1 LIMIT 1;",
		clientId,
	).Scan(&c.ID, &c.Name, &c.ClientId, &c.ClientSecret, &c.RedirectUri)
}

func (ac *AuthorizationCode) GetByCode(db *sql.DB, code string) error {
	var expiresAt string
	err := db.QueryRow(
		"SELECT id, code, expires_at, redirect_uri,client_id, user_account_id, active FROM authorization_codes WHERE code=$1 LIMIT 1;",
		code,
	).Scan(&ac.ID, &ac.Code, &expiresAt, &ac.RedirectUri, &ac.ClientId, &ac.UserAccountId, &ac.Active)
	if err != nil {
		return err
	}

	ac.expiresAt, err = parseDateIso(expiresAt)
	return err
}

func (ac *AuthorizationCode) Disable(db *sql.DB) error {
	_, err := db.Exec("UPDATE authorization_codes SET active=false WHERE id=$1;", ac.ID)
	return err
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

func (t *Token) GetByRefreshToken(db *sql.DB, refreshToken string) error {
	var accessTokenExpiresAt string
	var refreshTokenExpiresAt string
	err := db.QueryRow(
		"SELECT id, access_token, access_token_expires_at, refresh_token, refresh_token_expires_at, client_id, user_account_id, active FROM tokens WHERE refresh_token=$1 LIMIT 1;",
		refreshToken,
	).Scan(&t.ID, &t.AccessToken, &accessTokenExpiresAt, &t.RefreshToken, &refreshTokenExpiresAt, &t.ClientId, &t.UserAccountId, &t.Active)
	if err != nil {
		return err
	}
	t.AccessTokenExpiresAt, err = parseDateIso(accessTokenExpiresAt)
	if err != nil {
		return err
	}
	t.RefreshTokenExpiresAt, err = parseDateIso(refreshTokenExpiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (t *Token) GetByAccessToken(db *sql.DB, accessToken string) error {
	var accessTokenExpiresAt string
	var refreshTokenExpiresAt string
	err := db.QueryRow(
		"SELECT id, access_token, access_token_expires_at, refresh_token, refresh_token_expires_at, client_id, user_account_id, active FROM tokens WHERE access_token=$1 LIMIT 1;",
		accessToken,
	).Scan(&t.ID, &t.AccessToken, &accessTokenExpiresAt, &t.RefreshToken, &refreshTokenExpiresAt, &t.ClientId, &t.UserAccountId, &t.Active)
	if err != nil {
		return err
	}
	t.AccessTokenExpiresAt, err = parseDateIso(accessTokenExpiresAt)
	if err != nil {
		return err
	}
	t.RefreshTokenExpiresAt, err = parseDateIso(refreshTokenExpiresAt)
	if err != nil {
		return err
	}
	return nil
}

func nowPlusSeconds(seconds int) time.Time {
	return time.Now().Local().Add(time.Second * time.Duration(seconds))
}

func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func parseDateIso(date string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05-0700", date)
}

func formatDateIso(date time.Time) string {
	return date.Format("2006-01-02T15:04:05-0700")
}
