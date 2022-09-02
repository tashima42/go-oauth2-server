package handlers

import (
	"database/sql"
	"net/http"

	"github.com/tashima42/go-oauth2-server/data"
	"github.com/tashima42/go-oauth2-server/helpers"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler struct {
	DB *sql.DB
}

type LoginRequestDTO struct {
	Username        string `schema:"username"`
	Password        string `schema:"password"`
	Country         string `schema:"country"`
	RedirectUri     string `schema:"redirect_uri"`
	State           string `schema:"state"`
	ClientId        string `schema:"client_id"`
	ResponseType    string `schema:"response_type"`
	FailureRedirect string `schema:"failureRedirect"`
	CpConvert       string `schema:"cp_convert"`
}
type LoginResponseDTO struct {
	Success     bool   `json:"success"`
	RedirectUri string `json:"redirect_uri"`
	State       string `json:"state"`
	Code        string `json:"code"`
}

func (lh *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "LOGIN-PARSE-FORM-ERROR", err.Error())
		return
	}
	var loginRequest LoginRequestDTO
	helpers.Decoder.Decode(&loginRequest, r.PostForm)

	c := data.Client{ClientId: loginRequest.ClientId}
	err = c.GetByClientId(lh.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			helpers.RespondWithError(w, http.StatusBadRequest, "LOGIN-INVALID-CLIENT-ID", "Client Id not found")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "LOGIN-CLIENT-ID-FAILED", err.Error())
		}
		return
	}
	if c.RedirectUri != loginRequest.RedirectUri {
		helpers.RespondWithError(w, http.StatusBadRequest, "LOGIN-INVALID-REDIRECT-URI", "Invalid redirect_uri")
		return
	}

	u := data.UserAccount{Username: loginRequest.Username, Country: loginRequest.Country}
	err = u.GetByUsernameAndCountry(lh.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			helpers.RespondWithError(w, http.StatusBadRequest, "LOGIN-INVALID-USERNAME-OR-COUNTRY", "Username and country combination not found")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "LOGIN-USERNAME-OR-COUNTRY-FAILED", err.Error())
		}
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(loginRequest.Password)) != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "LOGIN-INVALID-PASSWORD", "invalid password")
		return
	}

	ac := data.AuthorizationCode{
		RedirectUri:   loginRequest.RedirectUri,
		ClientId:      c.ID,
		UserAccountId: u.ID,
		Code:          helpers.GenerateRandomString(64),
		ExpiresAt:     helpers.NowPlusSeconds(helpers.AuthorizationCodeExpiration),
	}
	err = ac.CreateAuthorizationCode(lh.DB)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "LOGIN-FAILED-CREATE-AUTHORIZATION-CODE", err.Error())
		return
	}

	loginResponse := LoginResponseDTO{Success: true, RedirectUri: ac.RedirectUri, State: loginRequest.State, Code: ac.Code}

	helpers.RespondWithJSON(w, http.StatusOK, loginResponse)
}

func (lh *LoginHandler) LoginCustom(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	u := data.UserAccount{Username: username}
	err := u.GetByUsername(lh.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			helpers.RespondWithError(w, http.StatusBadRequest, "LOGIN-INVALID-USERNAME-OR-COUNTRY", "Username and country combination not found")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "LOGIN-USERNAME-OR-COUNTRY-FAILED", err.Error())
		}
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "LOGIN-INVALID-PASSWORD", "invalid password")
		return
	}

	userInfoResponse := UserInfoResponseDTO{Success: true, SubscriberId: u.SubscriberId, CountryCode: u.Country}

	helpers.RespondWithJSON(w, http.StatusOK, userInfoResponse)
}
