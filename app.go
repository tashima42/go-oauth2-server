package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/tashima42/go-oauth2-server/handlers"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode='disable'", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	loginHandler := handlers.LoginHandler{DB: a.DB}
	tokenHandler := handlers.TokenHandler{DB: a.DB}
	userInfoHandler := handlers.UserInfoHandler{DB: a.DB}
	a.Router.HandleFunc("/auth/login", loginHandler.Login).Methods("POST")
	a.Router.HandleFunc("/auth/token", tokenHandler.Token).Methods("POST")
	a.Router.HandleFunc("/userinfo", userInfoHandler.UserInfo).Methods("GET")
	a.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/views/")))
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
