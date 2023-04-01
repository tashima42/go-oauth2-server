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
	// use differnt cors middleware for login, only accept same origin
	router.POST("/login", handler.Login)
	router.POST("/user-accounts", handler.CreateUserAccount)

	router.Use(handler.AuthMiddleware)
	router.POST(
		"/dev-accounts",
		handler.VerifyRequiredScopes([]string{helpers.DevAccountCreateScope}),
		handler.CreateDevAccount,
	)

	router.GET("/authorize", handler.Authorize)
	router.POST("/token", handler.Token)

	router.POST(
		"/clients",
		handler.VerifyRequiredScopes([]string{helpers.ClientCreateScope}),
		handler.CreateClient,
	)
	router.GET("/userinfo", handler.UserInfo)

	router.Run(":8096")
}
