package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tashima42/go-oauth2-server/db"
)

type CreateUserAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) CreateUserAccount(c *gin.Context) {
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
