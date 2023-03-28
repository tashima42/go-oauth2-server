package handlers

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

// TODO: uncomment after authorization middleware is implemented
// func (h *Handler) UserInfo(w http.ResponseWriter, r *http.Request) {
// 	accessToken := r.Header.Get("Authorization")
// 	splitToken := strings.Split(accessToken, "Bearer ")
// 	accessToken = splitToken[1]

// 	t := db.Token{AccessToken: accessToken}
// 	err := t.GetByAccessToken(uh.DB)
// 	if err != nil {
// 		helpers.RespondWithError(w, http.StatusInternalServerError, "USERINFO-INVALID-ACCESS-TOKEN", err.Error())
// 		return
// 	}

// 	u := db.UserAccount{ID: t.UserAccountID}
// 	err = u.GetById(uh.DB)
// 	if err != nil {
// 		helpers.RespondWithError(w, http.StatusInternalServerError, "USERINFO-FAILED-GET-USER", err.Error())
// 		return
// 	}

// 	userInfoResponse := UserInfoResponseDTO{Success: true, SubscriberId: u.SubscriberId, CountryCode: u.Country}
// 	helpers.RespondWithJSON(w, http.StatusOK, userInfoResponse)
// }
