package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tashima42/go-oauth2-server/db"
)

type UserInfoResponse struct {
	Username string `json:"username"`
}

func (h *Handler) UserInfo(c *gin.Context) {
	userRaw, exists := c.Get("userToken")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get user information"})
		return
	}
	user := userRaw.(db.UserAccount)
	userInfoResponse := UserInfoResponse{Username: user.Username}
	c.JSON(http.StatusOK, userInfoResponse)
}
