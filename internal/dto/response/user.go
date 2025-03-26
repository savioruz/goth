package response

type UserRegisterResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type UserLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
