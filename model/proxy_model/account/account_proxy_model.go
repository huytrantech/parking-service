package account

type LoginAccountRequest struct {
	Phone string `json:"phone"`
	Password string `json:"password"`
}

type LoginAccountResponse struct {
	Token string `json:"token"`
	Phone string `json:"phone"`
	Username string `json:"username"`
}
