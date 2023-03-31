package handlers

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestTokenAuthorizationCodeGrant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("SELECT id, name, client_id, client_secret, redirect_uri FROM clients").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
}
