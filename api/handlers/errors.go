package handlers

import "github.com/pkg/errors"

type Oauth2Error error

var ServerError Oauth2Error = errors.New("server_error")
var InvalidRequest Oauth2Error = errors.New("invalid_request")
var UnsupportedResponseType Oauth2Error = errors.New("unsupported_response_type")
var InvalidScope Oauth2Error = errors.New("invalid_scope")

var RedirectURIInvalid error = errors.New("redirect_uri is invalid")
var RedirectURIRequired error = errors.New("redirect_uri is required")
var ResponseTypeRequired error = errors.New("response_type is required")
var ClientIDRequired error = errors.New("client_id is required")
var GrantTypeRequired error = errors.New("grant_type is required")
var GrantTypeOneOf error = errors.New("grant_type must be one of 'authorization_code', 'refresh_token'")
var CodeRequired error = errors.New("code is required")
var RefreshTokenRequired error = errors.New("refresh_token is required")
