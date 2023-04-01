package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CORSMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
}

func (h *Handler) AuthMiddleware(c *gin.Context) {
	var accessToken string
	var err error

	accessToken = c.GetHeader("Authorization")
	if accessToken == "" {
		accessToken, err = c.Cookie("SESSION")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: missing access token"})
			return
		}
	}

	token, err := h.jwtHelper.VerifyToken(accessToken)
	if err != nil {
		log.Println("error verifying token: ", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid access token"})
		return
	}

	c.Set("token", token)
	c.Next()
}
