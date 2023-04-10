package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Login(c *gin.Context) {
	// TODO: fix all the rollbacks
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "errorCode": "LOGIN-PARSE-FORM-ERROR"})
		return
	}
	tx, err := h.repo.BeginTxx(c, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "LOGIN-TRANSACTION-ERROR"})
		return
	}

	userAccount, err := h.repo.GetUserAccountByUsernameTxx(tx, loginRequest.Username)
	if err != nil {
		if err = db.Rollback(tx, err); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "LOGIN-ROLLBACK-ERROR"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "errorCode": "LOGIN-INVALID-USERNAME"})
		return
	}

	if matches, err := h.hashHelper.Verify(loginRequest.Password, userAccount.Password); err != nil || !matches {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid password", "errorCode": "LOGIN-INVALID-PASSWORD"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "LOGIN-COMMIT-ERROR"})
		return
	}

	accessToken := db.Token{
		ClientID:    "",
		UserAccount: *userAccount,
		Scopes:      userAccount.ScopesToSlice(),
		ExpiresAt:   helpers.NowPlusSeconds(int(helpers.AccessTokenExpiration)),
	}
	accessTokenJWT, err := h.jwtHelper.GenerateToken(accessToken.ToMap())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "LOGIN-ACCESS-TOKEN-GENERATION-ERROR"})
		return
	}
	// TODO: check what happens if domain is not set
	c.SetCookie("SESSION", accessTokenJWT, int(helpers.AccessTokenExpiration), "/", "", true, true)
	c.JSON(http.StatusOK, LoginResponse{Token: accessTokenJWT})
}
