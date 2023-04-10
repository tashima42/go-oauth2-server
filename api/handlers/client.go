package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
)

type CreateClientRequest struct {
	ClientID    string `json:"clientID"`
	Name        string `json:"name"`
	RedirectURI string `json:"redirectURI"`
}
type CreateClientResponse struct {
	Name         string `json:"name"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	RedirectURI  string `json:"redirectURI"`
}

type ClientInfoRequest struct {
	ClientID string `json:"clientID"`
}
type ClientInfoResponse struct {
	ClientID string `json:"clientID"`
	Name     string `json:"name"`
}

func (h *Handler) CreateClient(c *gin.Context) {
	// TODO: fix all the rollbacks
	var createClientRequest CreateClientRequest
	if err := c.ShouldBindJSON(&createClientRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-CLIENT-PARSE-FORM-ERROR"})
		return
	}
	tx, err := h.repo.BeginTxx(c, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-CLIENT-TRANSACTION-ERROR"})
		return
	}

	existingClient, err := h.repo.GetClientByClientIDTxx(tx, createClientRequest.ClientID)
	if existingClient != nil {
		if err = db.Rollback(tx, err); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-CLIENT-ROLLBACK-ERROR"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "Client ID already exists", "errorCode": "CREATE-CLIENT-CLIENT-ID-ALREADY-EXISTS"})
		return
	}

	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			if err = db.Rollback(tx, err); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-CLIENT-ROLLBACK-ERROR"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-CLIENT-GET-CLIENT-ERROR"})
			return
		}
	}
	// TODO: validate redirect uri, make sure it's a valid url with https and it doesn't has a query string

	rawClientSecret, err := h.hashHelper.GenerateRandomString(128)
	if err != nil {
		if err = db.Rollback(tx, err); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-CLIENT-ROLLBACK-ERROR"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-CLIENT-FAILED-GENERATE-RANDOM-STRING"})
		return
	}
	hashedClientSecret, err := h.hashHelper.Hash(rawClientSecret)
	if err != nil {
		if err = db.Rollback(tx, err); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-CLIENT-ROLLBACK-ERROR"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-CLIENT-FAILED-HASH-CLIENT-SECRET"})
		return
	}

	client := db.Client{
		Name:         createClientRequest.Name,
		ClientID:     createClientRequest.ClientID,
		ClientSecret: hashedClientSecret,
		RedirectURI:  createClientRequest.RedirectURI,
		Scopes:       helpers.DefaultClientScopes,
	}

	err = h.repo.CreateClientTxx(tx, client)
	if err != nil {
		if err = db.Rollback(tx, err); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-CLIENT-ROLLBACK-ERROR"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-CLIENT-FAILED-TO-CREATE-CLIENT"})
		return
	}
	if err = tx.Commit(); err != nil {
		if err = db.Rollback(tx, err); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-CLIENT-ROLLBACK-ERROR"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-CLIENT-COMMIT-ERROR"})
		return
	}
	createClientResponse := CreateClientResponse{
		Name:         client.Name,
		ClientID:     client.ClientID,
		ClientSecret: rawClientSecret,
		RedirectURI:  client.RedirectURI,
	}
	c.JSON(http.StatusCreated, createClientResponse)
}

func (h *Handler) GetClientInfo(c *gin.Context) {
	clientInfoRequest := ClientInfoRequest{ClientID: c.Param("clientID")}
	client, err := h.repo.GetClientByClientID(c, clientInfoRequest.ClientID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "GET-CLIENT-INFO-GET-CLIENT-ERROR"})
		return
	}
	clientInfoResponse := ClientInfoResponse{
		ClientID: client.ClientID,
		Name:     client.Name,
	}
	c.JSON(http.StatusOK, clientInfoResponse)
}
