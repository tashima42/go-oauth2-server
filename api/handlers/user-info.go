package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tashima42/go-oauth2-server/db"
)

type UserInfoResponse struct {
	Username string   `json:"username"`
	Type     string   `json:"type"`
	Scopes   []string `json:"scopes"`
}

func (h *Handler) UserInfo(c *gin.Context) {
	rawToken, exists := c.Get("token")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: missing access token"})
		return
	}
	token := rawToken.(*db.Token)
	user := token.UserAccount
	userInfoResponse := UserInfoResponse{
		Username: user.Username,
		Type:     string(user.Type),
		Scopes:   token.Scopes,
	}
	c.JSON(http.StatusOK, userInfoResponse)
}
