package request

type UserRegisterRequest struct {
	Email    string `json:"email" validate:"required,email" example:"string@gmail.com"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"string@gmail.com"`
	Password string `json:"password" validate:"required,min=8"`
}
