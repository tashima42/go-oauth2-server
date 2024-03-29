package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
)

type TokenRequest struct {
	GrantType    string
	Code         string
	RedirectURI  string
	ClientID     string
	RefreshToken string
}
type TokenResponse struct {
	AccessToken           string             `json:"accessToken"`
	TokenType             string             `json:"tokenType"`
	ExpiresIn             helpers.Expiration `json:"expiresIn"`
	RefreshToken          string             `json:"refreshToken"`
	RefreshTokenExpiresIn helpers.Expiration `json:"refreshTokenExpiresIn"`
}

func (h *Handler) Token(c *gin.Context) {
	tokenRequest := TokenRequest{
		GrantType:    c.PostForm("grant_type"),
		Code:         c.PostForm("code"),
		RedirectURI:  c.PostForm("redirect_uri"),
		ClientID:     c.PostForm("client_id"),
		RefreshToken: c.PostForm("refresh_token"),
	}
	tx, err := h.repo.BeginTxx(c, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "TOKEN-TRANSACTION-ERROR"})
		return
	}
	rawClient, exists := c.Get("client")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: missing client", "errorCode": "TOKEN-INVALID-CLIENT"})
		return
	}
	client := rawClient.(*db.Client)

	log.Println("TokenRequest: ", tokenRequest)
	log.Println("ClientID: ", client.ClientID, " tokenRequest.ClientID: ", tokenRequest.ClientID)
	if client.ClientID != tokenRequest.ClientID {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: mismatched clientID in body and authorization header", "errorCode": "TOKEN-INVALID-CLIENT"})
		return
	}

	var userAccountID *string
	switch tokenRequest.GrantType {
	case string(helpers.AuthorizationCodeGrantType):
		userAccountID, err = h.authorizationCodeGrant(tx, tokenRequest, client)
		if err != nil || userAccountID == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "TOKEN-INVALID-AUTHORIZATION-CODE"})
			return
		}
	case string(helpers.RefreshTokenGrantType):
		userAccountID, err = h.refreshTokenGrant(tokenRequest)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "TOKEN-INVALID-REFRESH-TOKEN"})
			return
		}
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Grant Type", "errorCode": "TOKEN-INVALID-GRANT-TYPE"})
		return
	}

	userAccount, err := h.repo.GetUserAccountByIDTxx(tx, *userAccountID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "TOKEN-FAILED-TO-GET-USER-ACCOUNT"})
		return
	}

	accessToken := db.Token{
		ClientID:    client.ClientID,
		UserAccount: *userAccount,
		ExpiresAt:   helpers.NowPlusSeconds(int(helpers.AccessTokenExpiration)),
		Scopes:      client.Scopes,
	}
	refreshToken := db.Token{
		ClientID:    client.ClientID,
		UserAccount: *userAccount,
		ExpiresAt:   helpers.NowPlusSeconds(int(helpers.RefreshTokenExpiration)),
		Scopes:      client.Scopes,
	}
	accessTokenJWT, err := h.jwtHelper.GenerateToken(accessToken.ToMap())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "TOKEN-FAILED-TO-GENERATE-TOKEN"})
		return
	}
	refreshTokenJWT, err := h.jwtHelper.GenerateToken(refreshToken.ToMap())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "TOKEN-FAILED-TO-GENERATE-TOKEN"})
		return
	}

	tokenResponse := TokenResponse{
		TokenType:             "Bearer",
		AccessToken:           accessTokenJWT,
		ExpiresIn:             helpers.AccessTokenExpiration,
		RefreshToken:          refreshTokenJWT,
		RefreshTokenExpiresIn: helpers.RefreshTokenExpiration,
	}
	c.JSON(http.StatusOK, tokenResponse)
}

func (h *Handler) authorizationCodeGrant(tx *sqlx.Tx, tokenRequest TokenRequest, client *db.Client) (userAccountID *string, err error) {
	authorizationCode, err := h.repo.GetAuthorizationCodeByCodeAndClientTxx(tx, tokenRequest.Code, client.ID)
	if err != nil {
		return nil, err
	}
	if !authorizationCode.Active {
		return nil, errors.New("authorization code is not active")
	}
	err = h.repo.DisableAuthorizationCodeByIDTxx(tx, authorizationCode.ID)
	if err != nil {
		return nil, errors.New("failed to disable authorization code")
	}
	return &authorizationCode.UserAccountID, nil
}

func (h *Handler) refreshTokenGrant(tokenRequest TokenRequest) (userAccountID *string, err error) {
	token, err := h.jwtHelper.VerifyToken(tokenRequest.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &token.UserAccount.ID, nil
}

func (tr *TokenRequest) validate() error {
	if tr.GrantType == "" {
		return GrantTypeRequired
	}
	if tr.Code == "" {
		return CodeRequired
	}
	if tr.RedirectURI == "" {
		return RedirectURIRequired
	}
	if tr.ClientID == "" {
		return ClientIDRequired
	}
	if tr.GrantType != string(helpers.AuthorizationCodeGrantType) && tr.GrantType != string(helpers.RefreshTokenGrantType) {
		return GrantTypeOneOf
	}

	if tr.GrantType == string(helpers.RefreshTokenGrantType) && tr.RefreshToken == "" {
		return RefreshTokenRequired
	}
	return nil
}
