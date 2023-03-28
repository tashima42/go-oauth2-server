package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
)

type TokenRequest struct {
	ClientId     string `schema:"client_id"`
	ClientSecret string `schema:"client_secret"`
	GrantType    string `schema:"grant_type"`
	Code         string `schema:"code"`
	RefreshToken string `schema:"refresh_token"`
}
type TokenResponse struct {
	Success               bool               `json:"success"`
	TokenType             string             `json:"token_type"`
	AccessToken           string             `json:"access_token"`
	ExpiresIn             helpers.Expiration `json:"expires_in"`
	RefreshToken          string             `json:"refresh_token"`
	RefreshTokenExpiresIn helpers.Expiration `json:"refresh_token_expires_in"`
}

func (h *Handler) Token(c *gin.Context) {
	var tokenRequest TokenRequest
	if err := c.ShouldBindJSON(&tokenRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "errorCode": "TOKEN-PARSE-FORM-ERROR"})
	}
	tx, err := h.repo.BeginTxx(c, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "TOKEN-TRANSACTION-ERROR"})
	}
	client, err := h.repo.GetClientByClientIDTxx(tx, tokenRequest.ClientId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "TOKEN-FAILED-TO-GET-CLIENT"})
	}
	if matches, err := h.hashHelper.Verify(client.ClientSecret, tokenRequest.ClientSecret); err != nil || !matches {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Client Secret", "errorCode": "TOKEN-INVALID-CLIENT-SECRET"})
	}

	var userAccountID *string
	switch tokenRequest.GrantType {
	case string(helpers.AuthorizationCodeGrantType):
		userAccountID, err = h.authorizationCodeGrant(tx, tokenRequest)
		if err != nil || userAccountID == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "TOKEN-INVALID-AUTHORIZATION-CODE"})
		}
	case string(helpers.RefreshTokenGrantType):
		userAccountID, err = h.refreshTokenGrant(tokenRequest)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "TOKEN-INVALID-REFRESH-TOKEN"})
		}
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Grant Type", "errorCode": "TOKEN-INVALID-GRANT-TYPE"})
	}

	userAccount, err := h.repo.GetUserAccountByIDTxx(tx, *userAccountID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "TOKEN-FAILED-TO-GET-USER-ACCOUNT"})
	}

	accessToken := db.Token{
		ClientID:    client.ClientID,
		UserAccount: *userAccount,
		ExpiresAt:   helpers.NowPlusSeconds(int(helpers.AccessTokenExpiration)),
	}
	refreshToken := db.Token{
		ClientID:    client.ClientID,
		UserAccount: *userAccount,
		ExpiresAt:   helpers.NowPlusSeconds(int(helpers.RefreshTokenExpiration)),
	}
	accessTokenJWT, err := h.jwtHelper.GenerateToken(accessToken.ToMap())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "TOKEN-FAILED-TO-GENERATE-TOKEN"})
	}
	refreshTokenJWT, err := h.jwtHelper.GenerateToken(refreshToken.ToMap())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "TOKEN-FAILED-TO-GENERATE-TOKEN"})
	}

	tokenResponse := TokenResponse{
		Success:               true,
		TokenType:             "Bearer",
		AccessToken:           accessTokenJWT,
		ExpiresIn:             helpers.AccessTokenExpiration,
		RefreshToken:          refreshTokenJWT,
		RefreshTokenExpiresIn: helpers.RefreshTokenExpiration,
	}
	// TODO: check what happens if domain is not set
	c.SetCookie("SESSION", accessTokenJWT, int(helpers.AccessTokenExpiration), "/", "", true, true)
	c.JSON(http.StatusOK, tokenResponse)
}

func (h *Handler) authorizationCodeGrant(tx *sqlx.Tx, tokenRequest TokenRequest) (userAccountID *string, err error) {
	authorizationCode, err := h.repo.GetAuthorizationCodeByCodeTxx(tx, tokenRequest.Code)
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
