package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
)

type LoginRequest struct {
	Username        string `schema:"username"`
	Password        string `schema:"password"`
	Country         string `schema:"country"`
	RedirectURI     string `schema:"redirect_uri"`
	State           string `schema:"state"`
	ClientID        string `schema:"client_id"`
	ResponseType    string `schema:"response_type"`
	FailureRedirect string `schema:"failureRedirect"`
	CPConvert       string `schema:"cp_convert"`
}
type LoginResponse struct {
	RedirectURI string `json:"redirect_uri"`
	State       string `json:"state"`
	Code        string `json:"code"`
}

func (h *Handler) Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "errorCode": "LOGIN-PARSE-FORM-ERROR"})
	}
	if loginRequest.ResponseType != "code" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid response_type", "errorCode": "LOGIN-INVALID-RESPONSE-TYPE"})
	}
	tx, err := h.repo.BeginTxx(c, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "LOGIN-TRANSACTION-ERROR"})
	}
	client, err := h.repo.GetClientByClientIDTxx(tx, loginRequest.ClientID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "errorCode": "LOGIN-INVALID-CLIENT-ID"})
	}

	if client.RedirectURI != loginRequest.RedirectURI {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid redirect_uri", "errorCode": "LOGIN-INVALID-REDIRECT-URI"})
	}

	userAccount, err := h.repo.GetUserAccountByUsernameAndCountryTxx(tx, loginRequest.Username, loginRequest.Country)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "errorCode": "LOGIN-INVALID-USERNAME-OR-COUNTRY"})
	}

	if matches, err := h.hashHelper.Verify(userAccount.Password, loginRequest.Password); err != nil || !matches {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid password", "errorCode": "LOGIN-INVALID-PASSWORD"})
	}

	ac := db.AuthorizationCode{
		RedirectURI:   loginRequest.RedirectURI,
		ClientID:      client.ClientID,
		UserAccountID: userAccount.ID,
		Code:          helpers.GenerateRandomString(64),
		ExpiresAt:     helpers.NowPlusSeconds(int(helpers.AuthorizationCodeExpiration)),
	}
	err = h.repo.CreateAuthorizationCodeTxx(tx, ac)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "LOGIN-FAILED-CREATE-AUTHORIZATION-CODE"})
	}

	loginResponse := LoginResponse{RedirectURI: ac.RedirectURI, State: loginRequest.State, Code: ac.Code}
	c.JSON(http.StatusOK, loginResponse)
}
