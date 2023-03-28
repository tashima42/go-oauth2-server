package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tashima42/go-oauth2-server/db"
)

type UserInfoResponse struct {
	Success      bool   `json:"success"`
	SubscriberId string `json:"subscriber_id"`
	CountryCode  string `json:"country_code"`
}

func (h *Handler) UserInfo(c *gin.Context) {
	userRaw, exists := c.Get("userToken")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get user information"})
	}
	user := userRaw.(db.UserAccount)
	userInfoResponse := UserInfoResponse{
		Success:      true,
		SubscriberId: user.SubscriberID,
		CountryCode:  user.Country,
	}
	c.JSON(http.StatusOK, userInfoResponse)
}
