package jwt

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
)

type JWTHelper struct {
	secret []byte
}

func NewJWTHelperFromENV() (*JWTHelper, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}
	return &JWTHelper{secret: []byte(secret)}, nil
}

func (j *JWTHelper) GenerateToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTHelper) VerifyToken(tokenString string) (*db.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	expirationTime, err := claims.GetExpirationTime()
	if err != nil && expirationTime.Time.Unix() < time.Now().Unix() {
		return nil, errors.Wrap(err, "refresh token is expired")
	}
	log.Println("expirationTime", expirationTime.Time)
	parsedToken := db.Token{
		ExpiresAt:   expirationTime.Time,
		ClientID:    claims["clientID"].(string),
		Scopes:      claims["scopes"].([]helpers.Scope),
		UserAccount: db.UserAccountFromMap(claims["userAccount"].(map[string]interface{})),
	}
	return &parsedToken, nil
}
