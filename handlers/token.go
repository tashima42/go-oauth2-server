package handlers

import (
	"database/sql"
	"net/http"

	"github.com/tashima42/go-oauth2-server/data"
	"github.com/tashima42/go-oauth2-server/helpers"
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

	var c data.Client
	err = c.GetByClientId(th.DB, tokenRequest.ClientId)

	// add correct error validation with client_id not found message
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "TOKEN-INVALID-CLIENT-ID", err.Error())
		return
	}
	var userAccountId int
	if tokenRequest.GrantType == "authorization_code" {
		var ac data.AuthorizationCode
		err = ac.GetByCode(th.DB, tokenRequest.Code)
		// add correct error validation with client_id not found message
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, "TOKEN-INVALID-AUTHORIZATION-CODE", err.Error())
			return
		}
		err = ac.Disable(th.DB)
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, "TOKEN-FAILED-USE-AUTHORIZATION-CODE", err.Error())
			return
		}
		userAccountId = ac.UserAccountId
	} else if tokenRequest.GrantType == "refresh_token" {
		var t data.Token
		err = t.GetByRefreshToken(th.DB, tokenRequest.RefreshToken)
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, "TOKEN-INVALID-REFRESH-TOKEN", err.Error())
			return
		}
		// validate if refresh token is expired
		userAccountId = t.UserAccountId
	}
	token := data.Token{ClientId: c.ID, UserAccountId: userAccountId}
	err = token.CreateToken(th.DB)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "TOKEN-FAILED-TO-CREATE-TOKEN", err.Error())
		return
	}

	tokenResponse := TokenResponseDTO{
		Success:               true,
		TokenType:             "Bearer",
		AccessToken:           token.AccessToken,
		ExpiresIn:             86400,
		RefreshToken:          token.RefreshToken,
		RefreshTokenExpiresIn: 2628288,
	}
	helpers.RespondWithJSON(w, http.StatusOK, tokenResponse)
}
