package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
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

	// TODO: validate the client redirect_uri
	redirectLocation := authorizeRequest.RedirectURI

	// TODO: document the size of the code
	code, err := h.hashHelper.GenerateRandomString(64)
	if err != nil {
		redirectLocation = redirectLocation + fmt.Sprintf("?error=%s&error_description=%s", ServerError.Error(), "internal error while generating code")
		c.Redirect(http.StatusInternalServerError, url.QueryEscape(redirectLocation))
		return
	}

	redirectLocation = redirectLocation + fmt.Sprintf("?code=%s&state=%s", code, authorizeRequest.State)
	c.Redirect(http.StatusFound, url.QueryEscape(redirectLocation))
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
