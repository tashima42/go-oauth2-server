package handlers

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
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
	} else {
		// remove "Bearer " from the beginning of the token
		accessToken = accessToken[7:]
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

func (h *Handler) VerifyRequiredScopes(requiredScopes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawToken, exists := c.Get("token")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: missing access token"})
		}
		token := rawToken.(*db.Token)
		log.Println("requiredScopes: ", requiredScopes)
		log.Println("token.Scopes: ", token.Scopes)
		valid := helpers.SliceContainsSlice(requiredScopes, token.Scopes)
		log.Println("valid: ", valid)
		if !valid {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: missing required scope"})
			return
		}
		c.Next()
	}
}

func (h *Handler) BasicAuthClientMiddleware(c *gin.Context) {
	basicAuth := c.GetHeader("Authorization")
	if basicAuth == "" {
	} else {
		// remove "Basic " from the beginning of the token
		basicAuth = basicAuth[6:]
	}

	decodedAuth, err := base64.StdEncoding.DecodeString(basicAuth)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid client credentials"})
	}

	// split the decoded auth into clientID and clientSecret by the colon
	auth := strings.Split(string(decodedAuth), ":")
	clientID := auth[0]
	clientSecret := auth[1]

	client, err := h.repo.GetClientByClientID(c, clientID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid client credentials"})
	}

	if valid, err := h.hashHelper.Verify(clientSecret, client.ClientSecret); !valid || err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid client credentials"})
	}

	c.Set("client", client)
	c.Next()
}
