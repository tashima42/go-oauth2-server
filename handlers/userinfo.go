package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
)

type UserInfoHandler struct {
	DB *sql.DB
}

type UserInfoResponseDTO struct {
	Success      bool   `json:"success"`
	SubscriberId string `json:"subscriber_id"`
	CountryCode  string `json:"country_code"`
}

type AuthorizeResponseDTO struct {
	Access bool   `json:"access"`
	Rating string `json:"rating"`
	Ttl    int    `json:"ttl"`
}

type InvalidUserResponse struct {
	Access       bool   `json:"access"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func (uh *UserInfoHandler) UserInfo(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")
	splitToken := strings.Split(accessToken, "Bearer ")
	accessToken = splitToken[1]

	t := db.Token{AccessToken: accessToken}
	err := t.GetByAccessToken(uh.DB)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "USERINFO-INVALID-ACCESS-TOKEN", err.Error())
		return
	}

	u := db.UserAccount{ID: t.UserAccountId}
	err = u.GetById(uh.DB)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "USERINFO-FAILED-GET-USER", err.Error())
		return
	}

	userInfoResponse := UserInfoResponseDTO{Success: true, SubscriberId: u.SubscriberId, CountryCode: u.Country}
	helpers.RespondWithJSON(w, http.StatusOK, userInfoResponse)
}

func (uh *UserInfoHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	resourceId := r.URL.Query().Get("resource_id")
	subscriberId := r.URL.Query().Get("subscriber_id")

	if resourceId != "urn:tve:paramountplus" && resourceId != "urn:tve:starzbasic" {
		invalidUserResponse := InvalidUserResponse{Access: false, ErrorCode: "AUTHORIZATION-INVALID-RESOURCE-ID", ErrorMessage: "Invalid resource_id"}
		helpers.RespondWithJSON(w, http.StatusBadRequest, invalidUserResponse)
		return
	}

	u := db.UserAccount{SubscriberId: subscriberId}
	err := u.GetBySubscriberId(uh.DB)
	if err != nil {
		invalidUserResponse := InvalidUserResponse{Access: false, ErrorCode: "AUTHORIZATION-INVALID-SUBSCRIBER-ID", ErrorMessage: "Invalid subscriber_id"}
		helpers.RespondWithJSON(w, http.StatusBadRequest, invalidUserResponse)
		return
	}

	access := true

	if resourceId == "urn:tve:starzbasic" {
		access = false
	}

	authorizeResponse := AuthorizeResponseDTO{Access: access, Rating: "G", Ttl: 3600}
	helpers.RespondWithJSON(w, http.StatusOK, authorizeResponse)
}
