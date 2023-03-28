package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
	"golang.org/x/crypto/bcrypt"
)

type TokenHandler struct {
	DB *sql.DB
}

type TokenRequestDTO struct {
	ClientId     string `schema:"client_id"`
	ClientSecret string `schema:"client_secret"`
	GrantType    string `schema:"grant_type"`
	Code         string `schema:"code"`
	RefreshToken string `schema:"refresh_token"`
}
type TokenResponseDTO struct {
	Success               bool   `json:"success"`
	TokenType             string `json:"token_type"`
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
}

func (th *TokenHandler) Token(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "TOKEN-PARSE-FORM-ERROR", err.Error())
		return
	}
	var tokenRequest TokenRequestDTO
	helpers.Decoder.Decode(&tokenRequest, r.PostForm)

	c := db.Client{ClientID: tokenRequest.ClientId}
	err = c.GetByClientId(th.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			helpers.RespondWithError(w, http.StatusUnauthorized, "TOKEN-INVALID-CLIENT-ID", "Client id is invalid")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "TOKEN-INVALID-CLIENT-ID", err.Error())
		}
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(c.ClientSecret), []byte(tokenRequest.ClientSecret)) != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "TOKEN-INVALID-CLIENT-SECRET", "Invalid Client Secret")
		return
	}

	var userAccountId int
	if tokenRequest.GrantType == "authorization_code" {
		err = th.authorizationCodeGrant(tokenRequest, &userAccountId)
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, "TOKEN-INVALID-AUTHORIZATION-CODE", err.Error())
			return
		}
	} else if tokenRequest.GrantType == "refresh_token" {
		err = th.refreshTokenGrant(tokenRequest, &userAccountId)
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, "TOKEN-INVALID-REFRESH-TOKEN", err.Error())
			return
		}
	}
	token := db.Token{
		ClientId:              c.ID,
		UserAccountId:         userAccountId,
		AccessToken:           helpers.GenerateRandomString(64),
		RefreshToken:          helpers.GenerateRandomString(64),
		AccessTokenExpiresAt:  helpers.NowPlusSeconds(helpers.AccessTokenExpiration),
		RefreshTokenExpiresAt: helpers.NowPlusSeconds(helpers.RefreshTokenExpiration),
	}
	err = token.CreateToken(th.DB)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "TOKEN-FAILED-TO-CREATE-TOKEN", err.Error())
		return
	}

	tokenResponse := TokenResponseDTO{
		Success:               true,
		TokenType:             "Bearer",
		AccessToken:           token.AccessToken,
		ExpiresIn:             helpers.AccessTokenExpiration,
		RefreshToken:          token.RefreshToken,
		RefreshTokenExpiresIn: helpers.RefreshTokenExpiration,
	}
	helpers.RespondWithJSON(w, http.StatusOK, tokenResponse)
}

func (th *TokenHandler) authorizationCodeGrant(tokenRequest TokenRequestDTO, userAccountId *int) error {
	ac := db.AuthorizationCode{Code: tokenRequest.Code}
	err := ac.GetByCode(th.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return errors.New("authorization code not found")
		default:
			return errors.New("failed to get authorization code")
		}
	}
	if !ac.Active {
		return errors.New("authorization code is not active")
	}
	err = ac.Disable(th.DB)
	if err != nil {
		return errors.New("failed to disable authorization code")
	}
	*userAccountId = ac.UserAccountID
	return nil
}

func (th *TokenHandler) refreshTokenGrant(tokenRequest TokenRequestDTO, userAccountId *int) error {
	t := db.Token{RefreshToken: tokenRequest.RefreshToken}
	err := t.GetByRefreshToken(th.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return errors.New("authorization code not found")
		default:
			return errors.New("failed to get authorization code")
		}
	}
	if !t.Active {
		return errors.New("refresh token is not active")
	}
	err = t.Disable(th.DB)
	if err != nil {
		return errors.New("failed to disable refresh token")
	}
	*userAccountId = t.UserAccountId
	return nil
}
