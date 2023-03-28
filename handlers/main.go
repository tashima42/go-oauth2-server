package handlers

import "github.com/tashima42/go-oauth2-server/db"

type Handler struct {
	repo *db.Repo
}

func NewHandler(repo *db.Repo) *Handler {
	return &Handler{repo: repo}
}
