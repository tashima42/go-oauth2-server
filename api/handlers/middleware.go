package handlers

import (
	"log"
	"net/http"

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
		tokenScopesMap := make(map[string]bool)
		for _, scope := range token.Scopes {
			if scope == helpers.AdminScope {
				c.Next()
				return
			}
			tokenScopesMap[scope] = true
		}
		for _, requiredScope := range requiredScopes {
			if _, ok := tokenScopesMap[requiredScope]; !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: missing required scope"})
				return
			}
		}
		c.Next()
	}
}
