package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
)

type AuthorizeRequest struct {
	ResponseType string
	ClientID     string
	RedirectURI  string
	State        string
}

func (h *Handler) Authorize(c *gin.Context) {
	authorizeRequest := AuthorizeRequest{
		ResponseType: c.Query("response_type"),
		ClientID:     c.Query("client_id"),
		RedirectURI:  c.Query("redirect_uri"),
		State:        c.Query("state"),
	}
	err := authorizeRequest.validate()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tx, err := h.repo.BeginTxx(c, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	client, err := h.repo.GetClientByClientIDTxx(tx, authorizeRequest.ClientID)
	if err != nil {
		if err = db.Rollback(tx, err); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if authorizeRequest.RedirectURI == "" {
		authorizeRequest.RedirectURI = client.RedirectURI
	}
	redirectLocation := authorizeRequest.RedirectURI + "?"

	if authorizeRequest.RedirectURI != client.RedirectURI {
		redirectLocation = redirectLocation + url.QueryEscape(fmt.Sprintf("error=%s&error_description=%s", InvalidRequest.Error(), "redirect uri does not match"))
		c.Redirect(http.StatusBadRequest, redirectLocation)
		return
	}

	token, exists := c.Get("token")
	if !exists {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token not found"))
		return
	}
	parsedToken := token.(*db.Token)

	switch authorizeRequest.ResponseType {
	case "code":
		h.authorizeAuthorizationCodeGrant(c, authorizeRequest, redirectLocation, *parsedToken)
		return
	case "token":
		h.authorizeImplicitGrant(c, authorizeRequest, redirectLocation, *parsedToken)
		return
	}
	redirectLocation = redirectLocation + url.QueryEscape(fmt.Sprintf("error=%s&error_description=%s", UnsupportedResponseType.Error(), "response type not supported"))
	c.Redirect(http.StatusBadRequest, redirectLocation)
}

func (h *Handler) authorizeAuthorizationCodeGrant(c *gin.Context, authorizeRequest AuthorizeRequest, redirectLocation string, token db.Token) {
	// TODO: document the size of the code
	code, err := h.hashHelper.GenerateRandomString(64)
	if err != nil {
		redirectLocation = redirectLocation + url.QueryEscape(fmt.Sprintf("error=%s&error_description=%s", ServerError.Error(), "internal error while generating code"))
		c.Redirect(http.StatusInternalServerError, redirectLocation)
		return
	}

	redirectLocation = redirectLocation + url.QueryEscape(fmt.Sprintf("code=%s&state=%s", code, authorizeRequest.State))
	c.Redirect(http.StatusFound, redirectLocation)
}

func (h *Handler) authorizeImplicitGrant(c *gin.Context, authorizeRequest AuthorizeRequest, redirectLocation string, token db.Token) {
	// TODO: after implementing scopes, check if the client has the right scope
	accessToken := db.Token{
		ClientID:    authorizeRequest.ClientID,
		UserAccount: token.UserAccount,
		ExpiresAt:   helpers.NowPlusSeconds(int(helpers.AccessTokenExpiration)),
	}
	accessTokenJWT, err := h.jwtHelper.GenerateToken(accessToken.ToMap())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "errorCode": "TOKEN-FAILED-TO-GENERATE-TOKEN"})
		return
	}
	redirectLocation = strings.Replace(redirectLocation, "?", "#", 1) + url.QueryEscape(fmt.Sprintf("access_token=%s&token_type=bearer&expires_in=%d&state=%s", accessTokenJWT, helpers.AccessTokenExpiration, authorizeRequest.State))
	log.Println("redirectLocation", redirectLocation)
	c.Redirect(http.StatusFound, redirectLocation)
}

func (ar *AuthorizeRequest) validate() error {
	if ar.ResponseType == "" {
		return ResponseTypeRequired
	}
	if ar.ClientID == "" {
		return ClientIDRequired
	}
	if ar.RedirectURI != "" {
		if _, err := url.ParseRequestURI(ar.RedirectURI); err != nil {
			return RedirectURIInvalid
		}
	}
	return nil
}
