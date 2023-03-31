package helpers

type Expiration int
type GrantType string

const AccessTokenExpiration Expiration = 86400      // 1 day
const RefreshTokenExpiration Expiration = 2628288   // 1 month
const AuthorizationCodeExpiration Expiration = 3600 // 1 hour

const AuthorizationCodeGrantType GrantType = "authorization_code"
const RefreshTokenGrantType GrantType = "refresh_token"
