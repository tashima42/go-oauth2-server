package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
)

type AuthorizeRequest struct {
	ResponseType string
	ClientID     string
	RedirectURI  string
	State        string
	Scope        string
}

func (h *Handler) Authorize(c *gin.Context) {
	authorizeRequest := AuthorizeRequest{
		ResponseType: c.Query("response_type"),
		ClientID:     c.Query("client_id"),
		RedirectURI:  c.Query("redirect_uri"),
		State:        c.Query("state"),
		Scope:        c.Query("scope"),
	}
	log.Println("AuthorizeRequest: ", authorizeRequest)
	err := authorizeRequest.validate()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	requestedScopes := strings.Split(authorizeRequest.Scope, " ")

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
		c.Redirect(http.StatusFound, redirectLocation)
		return
	}

	validScopes := helpers.SliceContainsSlice(requestedScopes, client.Scopes)
	if !validScopes {
		if err = db.Rollback(tx, err); err != nil {
			redirectLocation = redirectLocation + url.QueryEscape(fmt.Sprintf("error=%s&error_description=%s", ServerError.Error(), "internal error while generating code"))
			c.Redirect(http.StatusFound, redirectLocation)
			return
		}
		redirectLocation = redirectLocation + url.QueryEscape(fmt.Sprintf("error=%s&error_description=%s", InvalidScope.Error(), "requested scopes are not a subset of client scopes"))
		log.Println("redirectLocation: ", redirectLocation)
		c.Redirect(http.StatusFound, redirectLocation)
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
		h.authorizeAuthorizationCodeGrant(c, tx, authorizeRequest, redirectLocation, *parsedToken)
		return
	case "token":
		h.authorizeImplicitGrant(c, tx, authorizeRequest, redirectLocation, *parsedToken, requestedScopes)
		return
	}
	redirectLocation = redirectLocation + url.QueryEscape(fmt.Sprintf("error=%s&error_description=%s", UnsupportedResponseType.Error(), "response type not supported"))
	c.Redirect(http.StatusFound, redirectLocation)
}

func (h *Handler) authorizeAuthorizationCodeGrant(c *gin.Context, tx *sqlx.Tx, authorizeRequest AuthorizeRequest, redirectLocation string, token db.Token) {
	// TODO: document the size of the code
	code, err := h.hashHelper.GenerateRandomString(64)
	if err != nil {
		redirectLocation = redirectLocation + url.QueryEscape(fmt.Sprintf("error=%s&error_description=%s", ServerError.Error(), "internal error while generating code"))
		c.Redirect(http.StatusFound, redirectLocation)
		return
	}

	client, err := h.repo.GetClientByClientIDTxx(tx, authorizeRequest.ClientID)
	if err != nil {
		db.Rollback(tx, err)
		redirectLocation = redirectLocation + url.QueryEscape(fmt.Sprintf("error=%s&error_description=%s", ServerError.Error(), "internal error while generating code"))
		c.Redirect(http.StatusFound, redirectLocation)
		return
	}

	expiresAt := helpers.NowPlusSeconds(int(helpers.AuthorizationCodeExpiration))
	authorizationCode := db.AuthorizationCode{
		Code:          code,
		ClientID:      client.ID,
		ExpiresAt:     pq.NullTime{Time: expiresAt, Valid: true},
		RedirectURI:   authorizeRequest.RedirectURI,
		UserAccountID: token.UserAccount.ID,
	}

	err = h.repo.CreateAuthorizationCodeTxx(tx, authorizationCode)
	if err != nil {
		db.Rollback(tx, err)
		redirectLocation = redirectLocation + url.QueryEscape(fmt.Sprintf("error=%s&error_description=%s", ServerError.Error(), "internal error while generating code"))
		c.Redirect(http.StatusFound, redirectLocation)
		return
	}

	err = tx.Commit()
	if err != nil {
		redirectLocation = redirectLocation + url.QueryEscape(fmt.Sprintf("error=%s&error_description=%s", ServerError.Error(), "internal error while generating code"))
		c.Redirect(http.StatusFound, redirectLocation)
		return
	}

	redirectLocation = redirectLocation + url.QueryEscape(fmt.Sprintf("code=%s&state=%s", code, authorizeRequest.State))
	c.Redirect(http.StatusFound, redirectLocation)
}

func (h *Handler) authorizeImplicitGrant(c *gin.Context, tx *sqlx.Tx, authorizeRequest AuthorizeRequest, redirectLocation string, token db.Token, scopes []string) {
	accessToken := db.Token{
		ClientID:    authorizeRequest.ClientID,
		UserAccount: token.UserAccount,
		ExpiresAt:   helpers.NowPlusSeconds(int(helpers.AccessTokenExpiration)),
		Scopes:      scopes,
	}
	accessTokenJWT, err := h.jwtHelper.GenerateToken(accessToken.ToMap())
	if err != nil {
		redirectLocation = redirectLocation + url.QueryEscape(fmt.Sprintf("error=%s&error_description=%s", ServerError.Error(), "internal error while generating token"))
		c.Redirect(http.StatusFound, redirectLocation)
	}

	err = tx.Commit()
	if err != nil {
		redirectLocation = redirectLocation + url.QueryEscape(fmt.Sprintf("error=%s&error_description=%s", ServerError.Error(), "internal error while generating code"))
		c.Redirect(http.StatusFound, redirectLocation)
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
