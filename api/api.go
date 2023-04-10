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
	api := router.Group("/api")
	api.Use(handler.CORSMiddleware)
	api.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })
	// use differnt cors middleware for login, only accept same origin
	api.POST("/login", handler.Login)
	api.POST("/user-accounts", handler.CreateUserAccount)

	api.Use(handler.AuthMiddleware)
	api.POST(
		"/dev-accounts",
		handler.VerifyRequiredScopes([]string{helpers.DevAccountCreateScope}),
		handler.CreateDevAccount,
	)

	api.GET("/clients/:clientID",
		handler.VerifyRequiredScopes([]string{helpers.ClientInfoReadScope}),
		handler.GetClientInfo,
	)

	api.GET("/authorize", handler.Authorize)
	api.POST("/token", handler.Token)

	api.POST(
		"/clients",
		handler.VerifyRequiredScopes([]string{helpers.ClientCreateScope}),
		handler.CreateClient,
	)
	api.GET(
		"/userinfo",
		handler.VerifyRequiredScopes([]string{helpers.UserAccountUserInfoReadScope}),
		handler.UserInfo,
	)

	router.Run(":8096")
}
