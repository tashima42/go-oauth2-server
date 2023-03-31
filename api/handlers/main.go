package handlers

import (
	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
	"github.com/tashima42/go-oauth2-server/helpers/jwt"
)

type Handler struct {
	repo       *db.Repo
	hashHelper *helpers.HashHelper
	jwtHelper  *jwt.JWTHelper
}

func NewHandler(repo *db.Repo, hashHelper *helpers.HashHelper, jwtHelper *jwt.JWTHelper) *Handler {
	return &Handler{repo: repo, hashHelper: hashHelper, jwtHelper: jwtHelper}
}
