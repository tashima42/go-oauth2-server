package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/tashima42/go-oauth2-server/data"
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

func (uh *UserInfoHandler) UserInfo(w http.ResponseWriter, r *http.Request) {

	accessToken := r.Header.Get("Authorization")
	splitToken := strings.Split(accessToken, "Bearer ")
	accessToken = splitToken[1]

	t := data.Token{}
	err := t.GetByAccessToken(uh.DB, accessToken)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "USERINFO-INVALID-ACCESS-TOKEN", err.Error())
		return
	}

	u := data.UserAccount{}
	err = u.GetById(uh.DB, t.UserAccountId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "USERINFO-FAILED-GET-USER", err.Error())
		return
	}

	userInfoResponse := UserInfoResponseDTO{Success: true, SubscriberId: u.SubscriberId, CountryCode: u.Country}
	helpers.RespondWithJSON(w, http.StatusOK, userInfoResponse)
}
