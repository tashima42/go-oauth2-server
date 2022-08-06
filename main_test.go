package main_test

import (
	//"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/tashima42/go-oauth2-server"
)

var a main.App

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func TestAuthorizePage(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestLogin(t *testing.T) {
	clearTable()

	password := "secret"
	u := main.UserAccount{Username: "user1@example.com", Password: password, Country: "AR", SubscriberId: "subscriber1"}
	u.CreateUserAccount(a.DB)
	c := main.Client{Name: "client name", ClientId: "client1", ClientSecret: "secret", RedirectUri: "https://tashima42.github.io/tbx-local-dummy/"}
	c.CreateClient(a.DB)

	state := "currentstate"

	data := url.Values{}
	data.Set("username", u.Username)
	data.Set("password", password)
	data.Set("country", u.Country)
	data.Set("redirect_uri", c.RedirectUri)
	data.Set("state", state)
	data.Set("client_id", c.ClientId)
	data.Set("failureRedirect", c.RedirectUri)
	data.Set("response_type", "code")
	data.Set("cp_convert", "dummy2")

	req, _ := http.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["success"] != true {
		t.Errorf("Expected 'success' to be true. Got %v", m["success"])
	}
	if m["redirect_uri"] != c.RedirectUri {
		t.Errorf("Expected 'redirect_uri' to be %v. Got %v", c.RedirectUri, m["redirect_uri"])
	}
	if m["state"] != state {
		t.Errorf("Expected 'state' to be %v. Got %v", state, m["state"])
	}
}

func TestCreateTokenWithAuthorizationCode(t *testing.T) {
	clearTable()

	u := main.UserAccount{Username: "user1@example.com", Password: "secret", Country: "AR", SubscriberId: "subscriber1"}
	u.CreateUserAccount(a.DB)
	c := main.Client{Name: "client name", ClientId: "client1", ClientSecret: "secret", RedirectUri: "https://tashima42.github.io/tbx-local-dummy/"}
	c.CreateClient(a.DB)
	ac := main.AuthorizationCode{ClientId: c.ID, RedirectUri: c.RedirectUri, UserAccountId: u.ID}
	ac.CreateAuthorizationCode(a.DB)

	data := url.Values{}
	data.Set("client_id", c.ClientId)
	data.Set("client_secret", c.ClientSecret)
	data.Set("code", ac.Code)
	data.Set("redirect_uri", ac.RedirectUri)
	data.Set("grant_type", "authorization_code")

	req, _ := http.NewRequest(http.MethodPost, "/auth/token", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["success"] != true {
		t.Errorf("Expected 'success' to be true. Got %v", m["success"])
	}
	if m["token_type"] != "Bearer" {
		t.Errorf("Expected 'token_type' to be Bearer. Got %v", m["token_type"])
	}
	if m["expires_in"] != float64(86400) {
		t.Errorf("Expected 'expires_in' to be 86400. Got %v", m["expires_in"])
	}
	if m["refresh_token_expires_in"] != float64(2628288) {
		t.Errorf("Expected 'refresh_token_expires_in' to be 2628288. Got %v", m["refresh_token_expires_in"])
	}
}

func TestCreateTokenWithRefreshToken(t *testing.T) {
	clearTable()

	u := main.UserAccount{Username: "user1@example.com", Password: "secret", Country: "AR", SubscriberId: "subscriber1"}
	u.CreateUserAccount(a.DB)
	c := main.Client{Name: "client name", ClientId: "client1", ClientSecret: "secret", RedirectUri: "https://tashima42.github.io/tbx-local-dummy/"}
	c.CreateClient(a.DB)
	tk := main.Token{ClientId: c.ID, UserAccountId: u.ID}
	tk.CreateToken(a.DB)

	data := url.Values{}
	data.Set("client_id", c.ClientId)
	data.Set("client_secret", c.ClientSecret)
	data.Set("refresh_token", tk.RefreshToken)
	data.Set("grant_type", "refresh_token")

	req, _ := http.NewRequest(http.MethodPost, "/auth/token", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["success"] != true {
		t.Errorf("Expected 'success' to be true. Got %v", m["success"])
	}
	if m["token_type"] != "Bearer" {
		t.Errorf("Expected 'token_type' to be Bearer. Got %v", m["token_type"])
	}
	if m["expires_in"] != float64(86400) {
		t.Errorf("Expected 'expires_in' to be 86400. Got %v", m["expires_in"])
	}
	if m["refresh_token_expires_in"] != float64(2628288) {
		t.Errorf("Expected 'refresh_token_expires_in' to be 2628288. Got %v", m["refresh_token_expires_in"])
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected code %d. Got %d\n", expected, actual)
	}
}

func ensureTableExists() {
	b, _ := ioutil.ReadFile("./schema.sql")
	tableCreationQuery := string(b)
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM authorization_codes;")
	a.DB.Exec("DELETE FROM tokens;")
	a.DB.Exec("DELETE FROM clients;")
	a.DB.Exec("DELETE FROM user_accounts;")
	a.DB.Exec("ALTER SEQUENCE authorization_codes_id_seq RESTART WITH 1")
	a.DB.Exec("ALTER SEQUENCE tokens_id_seq RESTART WITH 1")
	a.DB.Exec("ALTER SEQUENCE clients_id_seq RESTART WITH 1")
	a.DB.Exec("ALTER SEQUENCE user_accounts_id_seq RESTART WITH 1")
}
