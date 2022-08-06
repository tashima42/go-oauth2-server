package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode='disable'", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/auth/login", a.login).Methods("POST")
	a.Router.HandleFunc("/auth/token", a.token).Methods("POST")
	a.Router.HandleFunc("/userinfo", a.userInfo).Methods("GET")
	a.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/views/")))
}

var decoder = schema.NewDecoder()

func (a *App) login(w http.ResponseWriter, r *http.Request) {
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

	err := r.ParseForm()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "LOGIN-PARSE-FORM-ERROR", err.Error())
		return
	}
	var loginRequest LoginRequestDTO
	decoder.Decode(&loginRequest, r.PostForm)

	var c Client
	err = c.GetByClientId(a.DB, loginRequest.ClientId)
	// add correct error validation with client_id not found message
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "LOGIN-INVALID-CLIENT-ID", err.Error())
		return
	}
	if c.RedirectUri != loginRequest.RedirectUri {
		respondWithError(w, http.StatusBadRequest, "LOGIN-INVALID-REDIRECT-URI", "Invalid redirect_uri")
		return
	}

	var u UserAccount
	err = u.GetByUsernameAndCountry(a.DB, loginRequest.Username, loginRequest.Country)
	// add correct error validation with username with country not found message
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "LOGIN-INVALID-USERNAME-OR-COUNTRY", err.Error())
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(loginRequest.Password)) != nil {
		respondWithError(w, http.StatusUnauthorized, "LOGIN-INVALID-PASSWORD", "invalid password")
		return
	}

	ac := AuthorizationCode{RedirectUri: loginRequest.RedirectUri, ClientId: c.ID, UserAccountId: u.ID}
	err = ac.CreateAuthorizationCode(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "LOGIN-FAILED-CREATE-AUTHORIZATION-CODE", err.Error())
		return
	}

	loginResponse := LoginResponseDTO{Success: true, RedirectUri: ac.RedirectUri, State: loginRequest.State, Code: ac.Code}

	respondWithJSON(w, http.StatusOK, loginResponse)
}

func (a *App) token(w http.ResponseWriter, r *http.Request) {
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

	err := r.ParseForm()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "TOKEN-PARSE-FORM-ERROR", err.Error())
		return
	}
	var tokenRequest TokenRequestDTO
	decoder.Decode(&tokenRequest, r.PostForm)

	var c Client
	err = c.GetByClientId(a.DB, tokenRequest.ClientId)

	// add correct error validation with client_id not found message
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "TOKEN-INVALID-CLIENT-ID", err.Error())
		return
	}
	var userAccountId int
	if tokenRequest.GrantType == "authorization_code" {
		fmt.Printf("Authorization")
		var ac AuthorizationCode
		err = ac.GetByCode(a.DB, tokenRequest.Code)
		// add correct error validation with client_id not found message
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "TOKEN-INVALID-AUTHORIZATION-CODE", err.Error())
			return
		}
		err = ac.Disable(a.DB)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "TOKEN-FAILED-USE-AUTHORIZATION-CODE", err.Error())
			return
		}
		userAccountId = ac.UserAccountId
	} else if tokenRequest.GrantType == "refresh_token" {
		fmt.Printf("Refresh")
		var t Token
		err = t.GetByRefreshToken(a.DB, tokenRequest.RefreshToken)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "TOKEN-INVALID-REFRESH-TOKEN", err.Error())
			return
		}
		// validate if refresh token is expired
		userAccountId = t.UserAccountId
	}
	fmt.Printf("User: %v\n", userAccountId)
	token := Token{ClientId: c.ID, UserAccountId: userAccountId}
	err = token.CreateToken(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "TOKEN-FAILED-TO-CREATE-TOKEN", err.Error())
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
	respondWithJSON(w, http.StatusOK, tokenResponse)
}

func (a *App) userInfo(w http.ResponseWriter, r *http.Request) {
	type UserInfoResponseDTO struct {
		Success      bool   `json:"success"`
		SubscriberId string `json:"subscriber_id"`
		CountryCode  string `json:"country_code"`
	}

	accessToken := r.Header.Get("Authorization")
	splitToken := strings.Split(accessToken, "Bearer ")
	accessToken = splitToken[1]

	t := Token{}
	err := t.GetByAccessToken(a.DB, accessToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "USERINFO-INVALID-ACCESS-TOKEN", err.Error())
		return
	}

	u := UserAccount{}
	err = u.GetById(a.DB, t.UserAccountId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "USERINFO-FAILED-GET-USER", err.Error())
		return
	}

	userInfoResponse := UserInfoResponseDTO{Success: true, SubscriberId: u.SubscriberId, CountryCode: u.Country}
	respondWithJSON(w, http.StatusOK, userInfoResponse)
}

func respondWithError(w http.ResponseWriter, code int, errorCode string, message string) {
	fmt.Printf("ErrorCode: %v, Message: %v", errorCode, message)
	respondWithJSON(w, code, map[string]interface{}{"success": false, "errorCode": errorCode, "message": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
