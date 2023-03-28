package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/handlers"
	"github.com/tashima42/go-oauth2-server/helpers"
	"github.com/tashima42/go-oauth2-server/helpers/jwt"
)

func Serve(repo *db.Repo, hashHelper *helpers.HashHelper, jwtHelper *jwt.JWTHelper) {
	handler := handlers.NewHandler(repo, hashHelper, jwtHelper)
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(handler.CORSMiddleware)

	router.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })

	router.POST("/auth/login", handler.Login)
	router.POST("/auth/token", handler.Token)
	// // TODO: add user authorization middleware
	// router.GET("/userinfo", handler.UserInfo)

	// router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/views/")))
	router.Run(":8096")
}
