package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
)

type LoginRequest struct {
	Username     string `schema:"username"`
	Password     string `schema:"password"`
	RedirectURI  string `schema:"redirect_uri"`
	State        string `schema:"state"`
	ClientID     string `schema:"client_id"`
	ResponseType string `schema:"response_type"`
}
type LoginResponse struct {
	RedirectURI string `json:"redirectURI"`
	State       string `json:"state"`
	Code        string `json:"code"`
}

func (h *Handler) Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "errorCode": "LOGIN-PARSE-FORM-ERROR"})
		return
	}
	if loginRequest.ResponseType != "code" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid response_type", "errorCode": "LOGIN-INVALID-RESPONSE-TYPE"})
		return
	}
	tx, err := h.repo.BeginTxx(c, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "LOGIN-TRANSACTION-ERROR"})
		return
	}
	client, err := h.repo.GetClientByClientIDTxx(tx, loginRequest.ClientID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "errorCode": "LOGIN-INVALID-CLIENT-ID"})
		return
	}

	if client.RedirectURI != loginRequest.RedirectURI {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid redirect_uri", "errorCode": "LOGIN-INVALID-REDIRECT-URI"})
		return
	}

	userAccount, err := h.repo.GetUserAccountByUsernameTxx(tx, loginRequest.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "errorCode": "LOGIN-INVALID-USERNAME"})
		return
	}

	if matches, err := h.hashHelper.Verify(userAccount.Password, loginRequest.Password); err != nil || !matches {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid password", "errorCode": "LOGIN-INVALID-PASSWORD"})
		return
	}

	code, err := h.hashHelper.GenerateRandomString(64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "LOGIN-FAILED-GENERATE-RANDOM-STRING"})
		return
	}

	ac := db.AuthorizationCode{
		RedirectURI:   loginRequest.RedirectURI,
		ClientID:      client.ClientID,
		UserAccountID: userAccount.ID,
		Code:          code,
		ExpiresAt:     helpers.NowPlusSeconds(int(helpers.AuthorizationCodeExpiration)),
	}
	err = h.repo.CreateAuthorizationCodeTxx(tx, ac)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "LOGIN-FAILED-CREATE-AUTHORIZATION-CODE"})
		return
	}

	loginResponse := LoginResponse{RedirectURI: ac.RedirectURI, State: loginRequest.State, Code: ac.Code}
	c.JSON(http.StatusOK, loginResponse)
}
