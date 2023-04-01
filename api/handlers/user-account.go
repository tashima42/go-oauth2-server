package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
)

type CreateUserAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateDevAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) CreateUserAccount(c *gin.Context) {
	// TODO: fix all the rollbacks
	var createUserAccountRequest CreateUserAccountRequest
	if err := c.ShouldBindJSON(&createUserAccountRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-USER-ACCOUNT-PARSE-FORM-ERROR"})
		return
	}
	tx, err := h.repo.BeginTxx(c, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-USER-ACCOUNT-TRANSACTION-ERROR"})
		return
	}
	existingUser, err := h.repo.GetUserAccountByUsernameTxx(tx, createUserAccountRequest.Username)
	if existingUser != nil {
		if err = db.Rollback(tx, err); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-USER-ACCOUNT-ROLLBACK-ERROR"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "Username already exists", "errorCode": "CREATE-USER-ACCOUNT-USERNAME-ALREADY-EXISTS"})
		return
	}
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			if err = db.Rollback(tx, err); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-USER-ACCOUNT-ROLLBACK-ERROR"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-USER-ACCOUNT-GET-USER-ERROR"})
			return
		}
	}
	hashedPassword, err := h.hashHelper.Hash(createUserAccountRequest.Password)
	if err != nil {
		if err = db.Rollback(tx, err); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-USER-ACCOUNT-ROLLBACK-ERROR"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-USER-ACCOUNT-HASH-ERROR"})
		return
	}

	userAccount := db.UserAccount{
		Username: createUserAccountRequest.Username,
		Password: hashedPassword,
		Type:     db.UserAccountType,
		// TODO: review default scopes
		Scopes: []string{string(helpers.ClientCreateScope), string(helpers.ClientListScope)},
	}
	err = h.repo.CreateUserAccountTxx(tx, userAccount)
	if err != nil {
		if err = db.Rollback(tx, err); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-USER-ACCOUNT-ROLLBACK-ERROR"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-USER-ACCOUNT-CREATE-USER-ERROR"})
		return
	}
	if err = tx.Commit(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-USER-ACCOUNT-COMMIT-ERROR"})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) CreateDevAccount(c *gin.Context) {
	// TODO: fix all the rollbacks
	var createDevAccountRequest CreateDevAccountRequest
	if err := c.ShouldBindJSON(&createDevAccountRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-DEV-ACCOUNT-PARSE-FORM-ERROR"})
		return
	}
	tx, err := h.repo.BeginTxx(c, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-DEV-ACCOUNT-TRANSACTION-ERROR"})
		return
	}
	existingUser, err := h.repo.GetUserAccountByUsernameTxx(tx, createDevAccountRequest.Username)
	if existingUser != nil {
		if err = db.Rollback(tx, err); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-DEV-ACCOUNT-ROLLBACK-ERROR"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "Username already exists", "errorCode": "CREATE-DEV-ACCOUNT-USERNAME-ALREADY-EXISTS"})
		return
	}
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			if err = db.Rollback(tx, err); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-DEV-ACCOUNT-ROLLBACK-ERROR"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-DEV-ACCOUNT-GET-USER-ERROR"})
			return
		}
	}
	hashedPassword, err := h.hashHelper.Hash(createDevAccountRequest.Password)
	if err != nil {
		if err = db.Rollback(tx, err); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-DEV-ACCOUNT-ROLLBACK-ERROR"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-DEV-ACCOUNT-HASH-ERROR"})
		return
	}

	userAccount := db.UserAccount{
		Username: createDevAccountRequest.Username,
		Password: hashedPassword,
		Type:     db.DevAccountType,
		// TODO: review default scopes
		Scopes: []string{string(helpers.ClientCreateScope), string(helpers.ClientListScope)},
	}
	err = h.repo.CreateUserAccountTxx(tx, userAccount)
	if err != nil {
		if err = db.Rollback(tx, err); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-DEV-ACCOUNT-ROLLBACK-ERROR"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-DEV-ACCOUNT-CREATE-USER-ERROR"})
		return
	}
	if err = tx.Commit(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error(), "errorCode": "CREATE-DEV-ACCOUNT-COMMIT-ERROR"})
		return
	}
	c.Status(http.StatusCreated)
}
