package handlers

import (
	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
)

type Handler struct {
	repo       *db.Repo
	hashHelper *helpers.HashHelper
}

func NewHandler(repo *db.Repo, hashHelper *helpers.HashHelper) *Handler {
	return &Handler{repo: repo, hashHelper: hashHelper}
}
