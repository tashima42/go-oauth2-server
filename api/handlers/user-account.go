package handlers

type CreateUserAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
